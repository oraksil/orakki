package usecases

import (
	"errors"

	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
)

type SetupUseCase struct {
	ServiceConfig  *services.ServiceConfig
	MessageService services.MessageService
	WebRTCSession  services.WebRTCSession
	EngineFactory  engine.EngineFactory

	// localIceCandidatesQ chan
	gameId int64
}

func (uc *SetupUseCase) Prepare(prepare models.PrepareOrakki) (*models.Orakki, error) {
	uc.gameId = prepare.GameId

	return &models.Orakki{
		Id:    uc.ServiceConfig.MqRpcIdentifier,
		State: models.OrakkiStateReady,
	}, nil
}

func (uc *SetupUseCase) ProcessNewOffer(sdp models.SdpInfo) (*models.SdpInfo, error) {
	playerId := sdp.PeerId
	if playerId <= 0 {
		return nil, errors.New("invalid player id")
	}

	uc.WebRTCSession.StartIceProcess(
		playerId,
		uc.onLocalIceCandidate,
		uc.onIceConnectionStateChanged,
	)

	b64EncodedAnswer, err := uc.WebRTCSession.ProcessNewOffer(sdp)
	if err != nil {
		return nil, err
	}

	answerSdp := &models.SdpInfo{
		PeerId:           uc.gameId,
		SdpBase64Encoded: b64EncodedAnswer,
	}

	return answerSdp, nil
}

func (uc *SetupUseCase) ProcessRemoteIceCandidate(remoteIce models.IceCandidate) error {
	playerId := remoteIce.PeerId
	if playerId <= 0 {
		return errors.New("invalid player id")
	}

	return uc.WebRTCSession.ProcessRemoteIce(remoteIce)
}

func (uc *SetupUseCase) onLocalIceCandidate(b64EncodedIceCandidate string) {
	localIce := models.IceCandidate{
		PeerId:           uc.gameId,
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
