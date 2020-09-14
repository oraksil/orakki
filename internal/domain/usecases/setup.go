package usecases

import (
	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
)

type SetupUseCase struct {
	ServiceConfig  *services.ServiceConfig
	MessageService services.MessageService
	WebRTCSession  services.WebRTCSession
	EngineFactory  engine.EngineFactory
}

func (uc *SetupUseCase) ProcessNewOffer(sdp models.SdpInfo) (string, error) {
	uc.WebRTCSession.StartIceProcess(
		sdp.PeerId,
		uc.onLocalIceCandidate,
		uc.onIceConnectionStateChanged,
	)

	b64EncodedAnswer, err := uc.WebRTCSession.ProcessNewOffer(sdp)
	if err != nil {
		return "", err
	}

	return b64EncodedAnswer, nil
}

func (uc *SetupUseCase) ProcessRemoteIceCandidate(remoteIce models.IceCandidate) {
	uc.WebRTCSession.ProcessRemoteIce(remoteIce)
}

func (uc *SetupUseCase) onLocalIceCandidate(b64EncodedIceCandidate string) {
	localIce := models.IceCandidate{
		PeerId:           uc.ServiceConfig.PeerName,
		IceBase64Encoded: b64EncodedIceCandidate,
	}
	uc.MessageService.SendToAny(models.MSG_REMOTE_ICE_CANDIDATE, localIce)
}

func (uc *SetupUseCase) onIceConnectionStateChanged(connectionState string) {
	if connectionState == "connected" {
		rc := uc.WebRTCSession.GetRenderContext()
		ic := uc.WebRTCSession.GetInputContext()
		uc.EngineFactory.SetContexts(rc, ic)
	}
}
