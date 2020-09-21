package impl

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/oraksil/orakki/pkg/utils"
	"github.com/pion/webrtc/v2"
)

type PeerState int8

const (
	PEER_STATE_INIT             = 0
	PEER_STATE_TRACK_CONFIGURED = 1
)

type WebRTCSessionImpl struct {
	peers      map[int64]*webrtc.PeerConnection
	peerStates map[int64]PeerState

	videoTracks map[int]*webrtc.Track
	audioTracks map[int]*webrtc.Track

	peerInputEvents chan engine.InputEvent
}

func (s *WebRTCSessionImpl) StartIceProcess(peerId int64,
	onLocalIceCandidate func(b64EncodedIceCandidate string),
	onIceConnectionStateChanged func(connectionState string)) error {

	peer := s.getOrNewPeer(peerId)
	if peer == nil {
		return errors.New("peer connection does not exist.")
	}

	// Setup callback
	peer.OnICECandidate(func(c *webrtc.ICECandidate) {
		b64EncodedIceCandidate := ""
		if c != nil {
			b64EncodedIceCandidate, _ = utils.EncodeToB64EncodedJsonStr(c.ToJSON())
			fmt.Println("sending local ice candidate.")
		} else {
			fmt.Println("sending final ack for local ice candidate.")
		}
		onLocalIceCandidate(b64EncodedIceCandidate)
	})

	peer.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ice connection state changed: %s.\n", connectionState.String())
		onIceConnectionStateChanged(connectionState.String())
	})

	peer.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnOpen(func() {
			s.peerInputEvents <- engine.InputEvent{
				PlayerId: peerId,
				Type:     engine.InputEventTypeSessionOpen,
			}
		})

		d.OnClose(func() {
			s.peerInputEvents <- engine.InputEvent{
				PlayerId: peerId,
				Type:     engine.InputEventTypeSessionClose,
			}
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			s.peerInputEvents <- engine.InputEvent{
				PlayerId: peerId,
				Type:     engine.InputEventTypeKeyMessage,
				Data:     msg.Data,
			}
		})
	})

	return nil
}

func (s *WebRTCSessionImpl) ProcessNewOffer(sdp models.SdpInfo) (string, error) {
	peer := s.getOrNewPeer(sdp.PeerId)
	if peer == nil {
		return "", errors.New("peer connection does not exist.")
	}

	if peerState, _ := s.peerStates[sdp.PeerId]; peerState == PEER_STATE_TRACK_CONFIGURED {
		return "", errors.New("peer is already configured.")
	}

	offer := webrtc.SessionDescription{}
	err := utils.DecodeFromB64EncodedJsonStr(sdp.SdpBase64Encoded, &offer)
	if err != nil {
		return "", err
	}

	err = s.attachMatchedMediaTracksToPeer(peer, &offer)
	if err != nil {
		return "", errors.New("failed to extract and attach matched media tracks from peer.")
	}
	s.peerStates[sdp.PeerId] = PEER_STATE_TRACK_CONFIGURED

	err = peer.SetRemoteDescription(offer)
	if err != nil {
		return "", errors.New("failed to set remote decription")
	}

	answer, err := peer.CreateAnswer(nil)
	if err != nil {
		return "", errors.New("failed to create answer")
	}

	err = peer.SetLocalDescription(answer)
	if err != nil {
		return "", errors.New("failed to set local description")
	}

	// for firefox
	// answer.SDP = strings.ReplaceAll(answer.SDP, "a=sendrecv", "a=sendonly")

	b64EncodedAnswer, err := utils.EncodeToB64EncodedJsonStr(&answer)
	if err != nil {
		return "", err
	}

	fmt.Println("sdp offer is configured, returning sdp answer.")

	return b64EncodedAnswer, nil
}

func (s *WebRTCSessionImpl) ProcessRemoteIce(remoteIce models.IceCandidate) error {
	peer := s.getOrNewPeer(remoteIce.PeerId)
	if peer == nil {
		return errors.New("peer connection does not exist.")
	}

	var iceCandidate webrtc.ICECandidateInit
	err := utils.DecodeFromB64EncodedJsonStr(remoteIce.IceBase64Encoded, &iceCandidate)
	if err != nil {
		return err
	}

	err = peer.AddICECandidate(iceCandidate)
	if err != nil {
		return err
	}

	fmt.Println("player ice candidate added.")

	return nil
}

func (s *WebRTCSessionImpl) GetRenderContext() engine.RenderContext {
	audioTracks := make([]*webrtc.Track, 0, len(s.audioTracks))
	for _, track := range s.audioTracks {
		audioTracks = append(audioTracks, track)
	}

	videoTracks := make([]*webrtc.Track, 0, len(s.videoTracks))
	for _, track := range s.audioTracks {
		videoTracks = append(videoTracks, track)
	}

	return newWebRTCRenderContext(audioTracks, videoTracks)
}

func (s *WebRTCSessionImpl) GetInputContext() engine.InputContext {
	return newWebRTCInputContext(s.peerInputEvents)
}

func (s *WebRTCSessionImpl) getOrNewPeer(peerId int64) *webrtc.PeerConnection {
	peer, ok := s.peers[peerId]
	if ok {
		return peer
	}

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peer, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil
	}
	s.peers[peerId] = peer
	s.peerStates[peerId] = PEER_STATE_INIT

	return peer
}

func (s *WebRTCSessionImpl) attachMatchedMediaTracksToPeer(peer *webrtc.PeerConnection, offer *webrtc.SessionDescription) error {
	re := regexp.MustCompile(`a=rtpmap:([0-9]+) H264/90000`)
	matched := re.FindStringSubmatch(offer.SDP)
	payloadType, _ := strconv.ParseInt(matched[1], 10, 8)
	videoTrack := s.getOrNewVideoTrackByPayloadType(int(payloadType))

	re = regexp.MustCompile(`a=rtpmap:([0-9]+) opus/48000/2`)
	matched = re.FindStringSubmatch(offer.SDP)
	payloadType, _ = strconv.ParseInt(matched[1], 10, 8)
	audioTrack := s.getOrNewAudioTrackByPayloadType(int(payloadType))

	_, err := peer.AddTrack(videoTrack)
	if err != nil {
		return err
	}

	_, err = peer.AddTrack(audioTrack)
	if err != nil {
		return err
	}

	return nil
}

func (s *WebRTCSessionImpl) getOrNewVideoTrackByPayloadType(payloadType int) *webrtc.Track {
	if track, exists := s.videoTracks[payloadType]; exists {
		return track
	}

	codec := webrtc.NewRTPH264Codec(uint8(payloadType), 90000)
	track, _ := webrtc.NewTrack(uint8(payloadType), rand.Uint32(), "video", "pion2", codec)
	s.videoTracks[payloadType] = track

	return track
}

func (s *WebRTCSessionImpl) getOrNewAudioTrackByPayloadType(payloadType int) *webrtc.Track {
	if track, exists := s.audioTracks[payloadType]; exists {
		return track
	}

	codec := webrtc.NewRTPOpusCodec(uint8(payloadType), 8000)
	track, _ := webrtc.NewTrack(uint8(payloadType), rand.Uint32(), "audio", "pion2", codec)
	s.audioTracks[payloadType] = track

	return track
}

func NewWebRTCSession() services.WebRTCSession {
	return &WebRTCSessionImpl{
		peers:           make(map[int64]*webrtc.PeerConnection),
		peerStates:      make(map[int64]PeerState),
		audioTracks:     make(map[int]*webrtc.Track),
		videoTracks:     make(map[int]*webrtc.Track),
		peerInputEvents: make(chan engine.InputEvent, 1024),
	}
}
