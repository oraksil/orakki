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

	GameId int64
}

func (uc *SetupUseCase) Prepare(prepare models.PrepareOrakki) (*models.Orakki, error) {
	uc.GameId = prepare.GameId

	return &models.Orakki{
		Id:    uc.ServiceConfig.MqRpcIdentifier,
		State: models.OrakkiStateReady,
	}, nil
}

func (uc *SetupUseCase) ProcessNewOffer(sdp models.SdpInfo) (*models.SdpInfo, error) {
	uc.WebRTCSession.StartIceProcess(
		sdp.PeerId,
		uc.onLocalIceCandidate,
		uc.onIceConnectionStateChanged,
	)

	b64EncodedAnswer, err := uc.WebRTCSession.ProcessNewOffer(sdp)
	if err != nil {
		return nil, err
	}

	answerSdp := &models.SdpInfo{
		PeerId:           uc.GameId,
		SdpBase64Encoded: b64EncodedAnswer,
	}

	return answerSdp, nil
}

func (uc *SetupUseCase) ProcessRemoteIceCandidate(remoteIce models.IceCandidate) {
	uc.WebRTCSession.ProcessRemoteIce(remoteIce)
}

func (uc *SetupUseCase) onLocalIceCandidate(b64EncodedIceCandidate string) {
	localIce := models.IceCandidate{
		PeerId:           uc.GameId,
		IceBase64Encoded: b64EncodedIceCandidate,
	}
	uc.MessageService.SendToAny(models.MsgRemoteIceCandidate, localIce)
}

func (uc *SetupUseCase) onIceConnectionStateChanged(connectionState string) {
	if connectionState == "connected" {
		rc := uc.WebRTCSession.GetRenderContext()
		ic := uc.WebRTCSession.GetInputContext()
		uc.EngineFactory.SetContexts(rc, ic)
	}
}
