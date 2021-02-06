package handlers

import (
	"encoding/json"

	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/oraksil/orakki/internal/domain/usecases"
	"github.com/sangwonl/mqrpc"
)

type SetupHandler struct {
	ServiceConfig *services.ServiceConfig
	SetupUseCase  *usecases.SetupUseCase
}

func (h *SetupHandler) handlePrepareOrakki(ctx *mqrpc.Context) interface{} {
	var preparePayload models.PrepareOrakki
	json.Unmarshal(ctx.GetMessage().Payload, &preparePayload)

	state, _ := h.SetupUseCase.Prepare(preparePayload)

	return state
}

func (h *SetupHandler) handleSetupWithNewOffer(ctx *mqrpc.Context) interface{} {
	var sdpOffer models.SdpInfo
	json.Unmarshal(ctx.GetMessage().Payload, &sdpOffer)

	sdpAnswer, _ := h.SetupUseCase.ProcessSdpExchange(sdpOffer)

	return sdpAnswer
}

func (h *SetupHandler) handleRemoteIceCandidate(ctx *mqrpc.Context) interface{} {
	var remoteIce models.IceCandidate
	json.Unmarshal(ctx.GetMessage().Payload, &remoteIce)

	h.SetupUseCase.ProcessRemoteIceCandidate(remoteIce)

	return nil
}

func (h *SetupHandler) Routes() []Route {
	return []Route{
		{MsgType: models.MsgPrepareOrakki, Handler: h.handlePrepareOrakki},
		{MsgType: models.MsgSetupWithNewOffer, Handler: h.handleSetupWithNewOffer},
		{MsgType: models.MsgRemoteIceCandidate, Handler: h.handleRemoteIceCandidate},
	}
}
