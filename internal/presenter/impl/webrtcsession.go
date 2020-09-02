package impl

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/pion/webrtc/v2"
	"gitlab.com/oraksil/orakki/internal/domain/engine"
	"gitlab.com/oraksil/orakki/internal/domain/services"
	"gitlab.com/oraksil/orakki/pkg/utils"
)

type PeerState int8

const (
	PEER_STATE_INIT             = 0
	PEER_STATE_TRACK_CONFIGURED = 1
)

type WebRTCSessionImpl struct {
	peers      map[string]*webrtc.PeerConnection
	peerStates map[string]PeerState

	videoTracks map[int]*webrtc.Track
	audioTracks map[int]*webrtc.Track

	peerInputEvents chan engine.InputEvent
}

func (s *WebRTCSessionImpl) StartIceProcess(
	peerId string,
	onLocalIceCandidates func(peerId, b64EncodedIceCandidate string),
	onIceConnectionStateChanged func(peerId, connectionState string)) error {

	peer := s.getOrNewPeer(peerId)
	if peer == nil {
		return errors.New("peer connection does not exist.")
	}

	// Setup callback
	peer.OnICECandidate(func(c *webrtc.ICECandidate) {
		b64EncodedIceCandidate, _ := utils.EncodeB64EncodedJsonStr(c.ToJSON())
		onLocalIceCandidates(peerId, b64EncodedIceCandidate)
	})

	peer.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		onIceConnectionStateChanged(peerId, connectionState.String())
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

func (s *WebRTCSessionImpl) OnCreatedOffer(
	peerId string,
	b64EncodedOffer string,
	onCreatedAnswer func(peerId, b64EncodedAnswer string)) error {

	peer := s.getOrNewPeer(peerId)
	if peer == nil {
		return errors.New("peer connection does not exist.")
	}

	if peerState, _ := s.peerStates[peerId]; peerState == PEER_STATE_TRACK_CONFIGURED {
		return errors.New("peer is already configured.")
	}

	offer := webrtc.SessionDescription{}
	err := utils.DecodeB64EncodedJsonStr(b64EncodedOffer, &offer)
	if err != nil {
		return err
	}

	err = s.attachMatchedMediaTracksToPeer(peer, &offer)
	if err != nil {
		return errors.New("failed to extract and attach matched media tracks from peer.")
	}
	s.peerStates[peerId] = PEER_STATE_TRACK_CONFIGURED

	err = peer.SetRemoteDescription(offer)
	if err != nil {
		return errors.New("failed to set remote decription")
	}

	answer, err := peer.CreateAnswer(nil)
	if err != nil {
		return errors.New("failed to create answer")
	}

	// for firefox
	// answer.SDP = strings.ReplaceAll(answer.SDP, "a=sendrecv", "a=sendonly")

	b64EncodedAnswer, err := utils.EncodeB64EncodedJsonStr(&answer)
	if err != nil {
		return err
	}

	onCreatedAnswer(peerId, b64EncodedAnswer)

	return nil
}

func (s *WebRTCSessionImpl) OnRemoteIceCandidates(
	peerId string,
	b64EncodedIceCandidate string) error {

	peer := s.getOrNewPeer(peerId)
	if peer == nil {
		return errors.New("peer connection does not exist.")
	}

	var iceCandidate webrtc.ICECandidateInit
	err := utils.DecodeB64EncodedJsonStr(b64EncodedIceCandidate, &iceCandidate)
	if err != nil {
		return err
	}

	err = peer.AddICECandidate(iceCandidate)
	if err != nil {
		return err
	}

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

func (s *WebRTCSessionImpl) getOrNewPeer(peerId string) *webrtc.PeerConnection {
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
		peers:           make(map[string]*webrtc.PeerConnection),
		peerStates:      make(map[string]PeerState),
		audioTracks:     make(map[int]*webrtc.Track),
		videoTracks:     make(map[int]*webrtc.Track),
		peerInputEvents: make(chan engine.InputEvent, 1024),
	}
}
