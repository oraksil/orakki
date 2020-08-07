package handlers

import (
	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/orakki/internal/domain/models"
	"gitlab.com/oraksil/orakki/internal/domain/usecases"
	"gitlab.com/oraksil/orakki/internal/presenter/mq/dto"
)

type SystemHandler struct {
	SystemMonitorUseCase *usecases.SystemStateMonitorUseCase
}

func (h *SystemHandler) handleFethState(ctx *mqrpc.Context) interface{} {
	sysState, _ := h.SystemMonitorUseCase.GetSystemState()
	return dto.SystemStateToOrakkiState(sysState)
}

func (h *SystemHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_FETCH_ORAKKI_STATE, Handler: h.handleFethState},
	}
}
