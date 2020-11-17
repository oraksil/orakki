package impl

import (
	"fmt"

	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	"github.com/pkg/errors"
)

// WebRTCRenderContext
type WebRTCRenderContext struct {
	audioTracks []*webrtc.Track
	videoTracks []*webrtc.Track
}

func (rc *WebRTCRenderContext) WriteAudioFrame(buf []byte) error {
	for _, track := range rc.audioTracks {
		err := track.WriteSample(media.Sample{Data: buf, Samples: 960})
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *WebRTCRenderContext) WriteVideoFrame(buf []byte) error {
	for _, track := range rc.videoTracks {
		err := track.WriteSample(media.Sample{Data: buf, Samples: 1})
		if err != nil {
			return err
		}
	}

	return nil
}

func newWebRTCRenderContext(audioTracks, videoTracks []*webrtc.Track) engine.RenderContext {
	return &WebRTCRenderContext{
		audioTracks: audioTracks,
		videoTracks: videoTracks,
	}
}

// WebRTCInputContext
type WebRTCInputContext struct {
	inputEvents chan engine.InputEvent
}

func (ic *WebRTCInputContext) FetchInput() (engine.InputEvent, error) {
	input := <-ic.inputEvents
	return input, nil
}

func newWebRTCInputContext(inputEvents chan engine.InputEvent) engine.InputContext {
	return &WebRTCInputContext{
		inputEvents: inputEvents,
	}
}

// WebRTCSessionContext
type WebRTCSessionContext struct {
	peers      map[int64]*webrtc.PeerConnection
	peerStates map[int64]PeerState
}

func (sc *WebRTCSessionContext) CloseSession(playerId int64) error {
	peer, ok := sc.peers[playerId]
	if !ok {
		return errors.New(fmt.Sprintf("peer connection not found for %d", playerId))
	}

	return peer.Close()
}

func newWebRTCSessionContext(
	peers map[int64]*webrtc.PeerConnection, peerStates map[int64]PeerState) engine.SessionContext {

	return &WebRTCSessionContext{
		peers:      peers,
		peerStates: peerStates,
	}
}

// WebRTCFrontInterface (Facade of rendering, input handling, connection management)
type WebRTCFrontInterface struct {
	renderContext  engine.RenderContext
	inputContext   engine.InputContext
	sessionContext engine.SessionContext
}

func (f *WebRTCFrontInterface) WriteAudioFrame(buf []byte) error {
	// direct write to webrtc track
	return f.renderContext.WriteAudioFrame(buf)
}

func (f *WebRTCFrontInterface) WriteVideoFrame(buf []byte) error {
	return f.renderContext.WriteVideoFrame(buf)
}

func (f *WebRTCFrontInterface) FetchInput() (engine.InputEvent, error) {
	return f.inputContext.FetchInput()
}

func (f *WebRTCFrontInterface) CloseSession(playerId int64) error {
	return f.sessionContext.CloseSession(playerId)
}

func newWebRTCFrontInterface(
	rc engine.RenderContext, ic engine.InputContext, sc engine.SessionContext) engine.FrontInterface {

	return &WebRTCFrontInterface{
		renderContext:  rc,
		inputContext:   ic,
		sessionContext: sc,
	}
}

// WebRTCGameEngineFactory
type WebRTCGameEngineFactory struct {
	serviceConfig *services.ServiceConfig

	renderContext  engine.RenderContext
	inputContext   engine.InputContext
	sessionContext engine.SessionContext
}

func (f *WebRTCGameEngineFactory) SetContexts(
	rc engine.RenderContext, ic engine.InputContext, sc engine.SessionContext) {

	f.renderContext = rc
	f.inputContext = ic
	f.sessionContext = sc
}

func (f *WebRTCGameEngineFactory) CanCreateEngine() bool {
	return f.renderContext != nil && f.inputContext != nil && f.sessionContext != nil
}

func (f *WebRTCGameEngineFactory) CreateEngine() *engine.GameEngine {
	return engine.NewGameEngine(
		newWebRTCFrontInterface(f.renderContext, f.inputContext, f.sessionContext),
		NewGipanDriver(
			f.serviceConfig.GipanImageFramesIpcUri,
			f.serviceConfig.GipanSoundFramesIpcUri,
			f.serviceConfig.GipanCmdInputsIpcUri,
		),
	)
}

func NewGameEngineFactory(serviceConf *services.ServiceConfig) engine.EngineFactory {
	return &WebRTCGameEngineFactory{
		serviceConfig: serviceConf,
	}
}
