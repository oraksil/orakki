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

	orakkiPeerId int64
}

func (uc *SetupUseCase) Prepare(prepare models.PrepareOrakki) (*models.Orakki, error) {
	uc.orakkiPeerId = prepare.GameId

	return &models.Orakki{
		Id:    uc.ServiceConfig.MqRpcIdentifier,
		State: models.OrakkiStateReady,
	}, nil
}

func (uc *SetupUseCase) ProcessSdpExchange(sdp models.SdpInfo) (*models.SdpInfo, error) {
	if sdp.Peer.PlayerId <= 0 {
		return nil, errors.New("invalid player id")
	}

	uc.WebRTCSession.SetupIceHandlers(
		sdp.Peer,
		uc.onLocalIceCandidate,
		uc.onIceConnectionStateChanged,
	)

	answerSdp, err := uc.WebRTCSession.ProcessNewOffer(sdp)
	if err != nil {
		return nil, err
	}

	return answerSdp, nil
}

func (uc *SetupUseCase) ProcessRemoteIceCandidate(remoteIce models.IceCandidate) error {
	if remoteIce.Peer.PlayerId <= 0 {
		return errors.New("invalid player id")
	}

	return uc.WebRTCSession.ProcessRemoteIce(remoteIce)
}

func (uc *SetupUseCase) onLocalIceCandidate(iceCandidate models.IceCandidate) {
	uc.MessageService.SendToAny(models.MsgRemoteIceCandidate, iceCandidate)
}

func (uc *SetupUseCase) onIceConnectionStateChanged(peerInfo models.PeerInfo, connectionState string) {
	if connectionState == "connected" {
		rc := uc.WebRTCSession.GetRenderContext()
		ic := uc.WebRTCSession.GetInputContext()
		sc := uc.WebRTCSession.GetSessionContext()
		uc.EngineFactory.SetContexts(rc, ic, sc)
	}
}
