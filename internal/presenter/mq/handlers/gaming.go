package handlers

import (
	"encoding/json"

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
	var gameInfo models.GameInfo
	json.Unmarshal(ctx.GetMessage().Payload, &gameInfo)

	h.GamingUseCase.StartGame(&gameInfo)

	return nil
}

func (h *GamingHandler) Routes() []Route {
	return []Route{
		{MsgType: models.MsgStartGame, Handler: h.handleStartGame},
	}
}
