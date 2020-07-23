package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/orakki/internal/domain/models"
	"gitlab.com/oraksil/orakki/internal/domain/usecases"
)

type HelloHandler struct {
	GameCtrlUseCase *usecases.GameCtrlUseCase
}

func (h *HelloHandler) handleHello(ctx *mqrpc.Context) interface{} {
	msg := ctx.GetMessage()
	var value map[string]string
	json.Unmarshal(msg.Payload, &value)
	fmt.Println(value)

	return h.GameCtrlUseCase.Pong()
}

func (h *HelloHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_HELLO, Handler: h.handleHello},
	}
}
