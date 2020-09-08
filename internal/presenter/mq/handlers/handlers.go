package handlers

import (
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/usecases"
	"github.com/oraksil/orakki/internal/presenter/mq/dto"
	"github.com/sangwonl/mqrpc"
)

type SystemHandler struct {
	SystemMonitorUseCase *usecases.SystemStateMonitorUseCase
}

func (h *SystemHandler) handleFetchState(ctx *mqrpc.Context) interface{} {
	sysState, _ := h.SystemMonitorUseCase.GetSystemState()
	return dto.SystemStateToOrakkiState(sysState)
}

func (h *SystemHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_FETCH_ORAKKI_STATE, Handler: h.handleFetchState},
	}
}

type SetupHandler struct {
	SetupUseCase *usecases.SetupUseCase
}

func (h *SetupHandler) handleSetupOffer(ctx *mqrpc.Context) interface{} {
	return &dto.SetupAnswer{PeerId: "some id", Answer: "Some Answer"}
}

func (h *SetupHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_HANDLE_SETUP_OFFER, Handler: h.handleSetupOffer},
	}
}
