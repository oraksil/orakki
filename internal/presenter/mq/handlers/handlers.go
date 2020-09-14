package handlers

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/oraksil/orakki/internal/domain/usecases"
	"github.com/sangwonl/mqrpc"
)

type SystemHandler struct {
	SystemUseCase *usecases.SystemUseCase
}

func (h *SystemHandler) handleFetchState(ctx *mqrpc.Context) interface{} {
	state, _ := h.SystemUseCase.GetOrakkiState()
	return state
}

func (h *SystemHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_FETCH_ORAKKI_STATE, Handler: h.handleFetchState},
	}
}

type SetupHandler struct {
	ServiceConfig *services.ServiceConfig
	SetupUseCase  *usecases.SetupUseCase
}

func (h *SetupHandler) handleSetupWithNewOffer(ctx *mqrpc.Context) interface{} {
	var newOffer map[string]string
	json.Unmarshal(ctx.GetMessage().Payload, &newOffer)

	var sdpOffer models.SdpInfo
	mapstructure.Decode(newOffer, &sdpOffer)

	b64EncodedAnswer, err := h.SetupUseCase.ProcessNewOffer(sdpOffer)

	if err == nil && len(b64EncodedAnswer) > 0 {
		return &models.SdpInfo{
			PeerId:           h.ServiceConfig.PeerName,
			SdpBase64Encoded: b64EncodedAnswer,
		}
	}

	return nil
}

func (h *SetupHandler) handleRemoteIceCandidate(ctx *mqrpc.Context) interface{} {
	// var newRemoteIce map[string]string
	// json.Unmarshal(ctx.GetMessage().Payload, &newRemoteIce)

	var remoteIce models.IceCandidate
	mapstructure.Decode(ctx.GetMessage().Payload, &remoteIce)

	h.SetupUseCase.ProcessRemoteIceCandidate(remoteIce)

	return nil
}

func (h *SetupHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_SETUP_WITH_NEW_OFFER, Handler: h.handleSetupWithNewOffer},
		{MsgType: models.MSG_REMOTE_ICE_CANDIDATE, Handler: h.handleRemoteIceCandidate},
	}
}
