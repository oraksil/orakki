package handlers

import (
	"encoding/json"
	"fmt"

	"gitlab.com/oraksil/orakki/internal/domain/models"
	"gitlab.com/oraksil/orakki/internal/domain/usecases"
	"gitlab.com/oraksil/sil/backend/pkg/mq"
)

type HelloHandler struct {
	GameCtrlUseCase *usecases.GameCtrlUseCase
}

func (h *HelloHandler) handleHello(ctx *mq.Context) interface{} {
	msg := ctx.GetMessage()
	var value map[string]string
	json.Unmarshal(msg.Payload, &value)
	fmt.Println(value)

	return h.GameCtrlUseCase.Pong()
}

func (h *HelloHandler) Routes() []mq.Route {
	return []mq.Route{
		{MsgType: models.MSG_HELLO, Handler: h.handleHello},
	}
}
