package impl

import (
	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
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

// WebRTCRenderer
type WebRTCRenderer struct {
	renderContext engine.RenderContext
}

func (r *WebRTCRenderer) WriteAudioFrame(buf []byte) error {
	// direct write to webrtc track
	return r.renderContext.WriteAudioFrame(buf)
}

func (r *WebRTCRenderer) WriteVideoFrame(buf []byte) error {
	return r.renderContext.WriteVideoFrame(buf)
}

func newWebRTCRenderer(rc engine.RenderContext) engine.Renderer {
	return &WebRTCRenderer{
		renderContext: rc,
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

// WebRTCInputHandler
type WebRTCInputHandler struct {
	inputContext engine.InputContext
}

func (i *WebRTCInputHandler) FetchInput() (engine.InputEvent, error) {
	return i.inputContext.FetchInput()
}

func newWebRTCInputHandler(inputContext engine.InputContext) engine.InputHandler {
	return &WebRTCInputHandler{
		inputContext: inputContext,
	}
}

// WebRTCGameEngineFactory
type WebRTCGameEngineFactory struct {
	serviceConfig *services.ServiceConfig

	renderContext engine.RenderContext
	inputContext  engine.InputContext
}

func (f *WebRTCGameEngineFactory) SetContexts(rc engine.RenderContext, ic engine.InputContext) {
	f.renderContext = rc
	f.inputContext = ic
}

func (f *WebRTCGameEngineFactory) CanCreateEngine() bool {
	return f.renderContext != nil && f.inputContext != nil
}

func (f *WebRTCGameEngineFactory) CreateEngine() *engine.GameEngine {
	return engine.NewGameEngine(
		newWebRTCRenderer(f.renderContext),
		newWebRTCInputHandler(f.inputContext),
		NewGipanDriver(
			f.serviceConfig.GipanImageFramesIpcUri,
			f.serviceConfig.GipanSoundFramesIpcUri,
			f.serviceConfig.GipanKeyInputsIpcUri,
		),
	)
}

func NewGameEngineFactory(serviceConf *services.ServiceConfig) engine.EngineFactory {
	return &WebRTCGameEngineFactory{
		serviceConfig: serviceConf,
	}
}
