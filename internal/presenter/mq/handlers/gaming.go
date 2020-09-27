package handlers

import (
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/oraksil/orakki/internal/domain/usecases"
	"github.com/sangwonl/mqrpc"
)

type GamingHandler struct {
	ServiceConfig *services.ServiceConfig
	GamingUseCase *usecases.GamingUseCase
}

func (h *GamingHandler) handleStartGame(ctx *mqrpc.Context) interface{} {
	h.GamingUseCase.StartGame()

	return nil
}

func (h *GamingHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MsgStartGame, Handler: h.handleStartGame},
	}
}
