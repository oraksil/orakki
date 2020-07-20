package handlers

import (
	"gitlab.com/oraksil/orakki/internal/domain/models"
	"gitlab.com/oraksil/orakki/internal/domain/usecases"
	"gitlab.com/oraksil/sil/backend/pkg/mq"
)

type HelloHandler struct {
	GameCtrlUseCase *usecases.GameCtrlUseCase
}

func (h *HelloHandler) handleHello(ctx *mq.Context) {
	msg := ctx.GetMessage()
	h.GameCtrlUseCase.Pong(&msg)
}

func (h *HelloHandler) Routes() []mq.Route {
	return []mq.Route{
		{MsgType: models.MSG_HELLO, Handler: h.handleHello},
	}
}
