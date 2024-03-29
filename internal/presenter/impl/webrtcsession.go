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
	PEER_STATE_INIT            = 0
	PEER_STATE_SDP_CONFIGURING = 1
	PEER_STATE_SDP_CONFIGURED  = 2
)

type WebRTCSessionImpl struct {
	peers      map[int64]*webrtc.PeerConnection
	peerStates map[int64]PeerState

	videoTracks map[int]*webrtc.Track
	audioTracks map[int]*webrtc.Track

	peerInputEvents chan engine.InputEvent

	turnUri       string
	turnSecretKey string
	turnTTL       int
}

func (s *WebRTCSessionImpl) SetupIceHandlers(
	peerInfo models.PeerInfo,
	onLocalIceCandidate func(iceCandidate models.IceCandidate),
	onIceConnectionStateChanged func(peerInfo models.PeerInfo, connectionState string)) error {

	peer := s.getOrNewPeer(peerInfo.PlayerId, true)
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
		onLocalIceCandidate(models.IceCandidate{
			Peer:             peerInfo,
			IceBase64Encoded: b64EncodedIceCandidate,
		})
	})

	peer.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ice connection state changed: %s.\n", connectionState.String())
		onIceConnectionStateChanged(peerInfo, connectionState.String())
	})

	peer.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnOpen(func() {
			s.peerInputEvents <- engine.InputEvent{
				PlayerId: peerInfo.PlayerId,
				Type:     engine.InputEventTypeSessionOpen,
			}
		})

		d.OnClose(func() {
			s.peerInputEvents <- engine.InputEvent{
				PlayerId: peerInfo.PlayerId,
				Type:     engine.InputEventTypeSessionClose,
			}
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			keyword := string(msg.Data[:])
			switch keyword {
			case "ping":
				s.peerInputEvents <- engine.InputEvent{
					PlayerId: peerInfo.PlayerId,
					Type:     engine.InputEventTypeHealthCheck,
					Data:     msg.Data,
				}
			default:
				s.peerInputEvents <- engine.InputEvent{
					PlayerId: peerInfo.PlayerId,
					Type:     engine.InputEventTypeKeyMessage,
					Data:     msg.Data,
				}
			}
		})
	})

	return nil
}

func (s *WebRTCSessionImpl) ProcessNewOffer(sdp models.SdpInfo) (*models.SdpInfo, error) {
	resetPeerState := func() {
		s.peerStates[sdp.Peer.PlayerId] = PEER_STATE_INIT
	}

	peer := s.getOrNewPeer(sdp.Peer.PlayerId, false)
	if peer == nil {
		return nil, errors.New("peer connection does not exist.")
	}

	if peerState, _ := s.peerStates[sdp.Peer.PlayerId]; peerState != PEER_STATE_INIT {
		return nil, errors.New("peer is already configured or being configured.")
	}
	s.peerStates[sdp.Peer.PlayerId] = PEER_STATE_SDP_CONFIGURING

	offer := webrtc.SessionDescription{}
	err := utils.DecodeFromB64EncodedJsonStr(sdp.SdpBase64Encoded, &offer)
	if err != nil {
		resetPeerState()
		return nil, err
	}

	err = s.attachMatchedMediaTracksToPeer(peer, &offer)
	if err != nil {
		resetPeerState()
		return nil, errors.New("failed to extract and attach matched media tracks from peer.")
	}

	err = peer.SetRemoteDescription(offer)
	if err != nil {
		resetPeerState()
		return nil, errors.New("failed to set remote decription")
	}

	answer, err := peer.CreateAnswer(nil)
	if err != nil {
		resetPeerState()
		return nil, errors.New("failed to create answer")
	}

	err = peer.SetLocalDescription(answer)
	if err != nil {
		resetPeerState()
		return nil, errors.New("failed to set local description")
	}

	// for firefox
	// answer.SDP = strings.ReplaceAll(answer.SDP, "a=sendrecv", "a=sendonly")

	b64EncodedAnswer, err := utils.EncodeToB64EncodedJsonStr(&answer)
	if err != nil {
		resetPeerState()
		return nil, err
	}

	fmt.Println("sdp offer is configured, returning sdp answer.")
	s.peerStates[sdp.Peer.PlayerId] = PEER_STATE_SDP_CONFIGURED

	return &models.SdpInfo{
		Peer:             sdp.Peer,
		SdpBase64Encoded: b64EncodedAnswer,
	}, nil
}

func (s *WebRTCSessionImpl) ProcessRemoteIce(remoteIce models.IceCandidate) error {
	peer := s.getOrNewPeer(remoteIce.Peer.PlayerId, false)
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
	for _, track := range s.videoTracks {
		videoTracks = append(videoTracks, track)
	}

	return newWebRTCRenderContext(audioTracks, videoTracks)
}

func (s *WebRTCSessionImpl) GetInputContext() engine.InputContext {
	return newWebRTCInputContext(s.peerInputEvents)
}

func (s *WebRTCSessionImpl) GetSessionContext() engine.SessionContext {
	return newWebRTCSessionContext(s.peers, s.peerStates)
}

func (s *WebRTCSessionImpl) releasePeer(peer *webrtc.PeerConnection, peerId int64) {
	peer.Close()

	delete(s.peers, peerId)
	delete(s.peerStates, peerId)
}

func (s *WebRTCSessionImpl) getOrNewPeer(peerId int64, new bool) *webrtc.PeerConnection {
	peer, ok := s.peers[peerId]
	if ok {
		if !new {
			return peer
		}

		s.releasePeer(peer, peerId)
	}

	iceServers := []webrtc.ICEServer{
		// {
		// URLs: []string{"stun:stun.l.google.com:19302"},
		// },
	}

	if s.turnUri != "" {
		username, password := utils.NewTurnAuth(
			strconv.FormatInt(peerId, 10),
			s.turnSecretKey,
			s.turnTTL,
		)

		iceServers = append(iceServers, webrtc.ICEServer{
			URLs:       []string{s.turnUri},
			Username:   username,
			Credential: password,
		})
	}

	config := webrtc.Configuration{ICEServers: iceServers}
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
	if matched == nil {
		return errors.New("no matched video track info from offer.")
	}

	payloadType, _ := strconv.ParseInt(matched[1], 10, 8)
	videoTrack := s.getOrNewVideoTrackByPayloadType(int(payloadType))
	_, err := peer.AddTrack(videoTrack)
	if err != nil {
		return err
	}

	re = regexp.MustCompile(`a=rtpmap:([0-9]+) opus/48000/2`)
	matched = re.FindStringSubmatch(offer.SDP)
	if matched == nil {
		return errors.New("no matched audio track info from offer.")
	}

	payloadType, _ = strconv.ParseInt(matched[1], 10, 8)
	audioTrack := s.getOrNewAudioTrackByPayloadType(int(payloadType))
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

func NewWebRTCSession(turnUri, turnSecretKey string, turnTTL int) services.WebRTCSession {
	return &WebRTCSessionImpl{
		peers:           make(map[int64]*webrtc.PeerConnection),
		peerStates:      make(map[int64]PeerState),
		audioTracks:     make(map[int]*webrtc.Track),
		videoTracks:     make(map[int]*webrtc.Track),
		peerInputEvents: make(chan engine.InputEvent, 1024),
		turnUri:         turnUri,
		turnSecretKey:   turnSecretKey,
		turnTTL:         turnTTL,
	}
}
