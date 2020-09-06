package usecases

import (
	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/services"
)

type SetupUseCase struct {
	// ServiceConfig  *services.ServiceConfig
	// MessageService services.MessageService
	WebRTCSession services.WebRTCSession
	EngineFactory engine.EngineFactory
}

func (uc *SetupUseCase) OnCreatedOffer(peerId, b64EncodedOffer string) {
	uc.WebRTCSession.StartIceProcess(
		peerId,
		uc.onLocalIceCandidates,
		uc.onIceConnectionStateChanged,
	)

	uc.WebRTCSession.OnCreatedOffer(
		peerId,
		b64EncodedOffer,
		uc.onCreatedAnswer,
	)
}

func (uc *SetupUseCase) OnRemoteIceCandidates(peerId, b64EncodedIceCandidate string) {
	uc.WebRTCSession.OnRemoteIceCandidates(
		peerId,
		b64EncodedIceCandidate,
	)
}

func (uc *SetupUseCase) onCreatedAnswer(peerId, b64EncodedAnswer string) {
}

func (uc *SetupUseCase) onLocalIceCandidates(peerId, b64EncodedIceCandidate string) {
}

func (uc *SetupUseCase) onIceConnectionStateChanged(peerId, connectionState string) {
	if connectionState == "connected" {
		rc := uc.WebRTCSession.GetRenderContext()
		ic := uc.WebRTCSession.GetInputContext()
		uc.EngineFactory.SetContexts(rc, ic)
	}
}
