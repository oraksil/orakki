package usecases

import (
	"encoding/json"

	"gitlab.com/oraksil/sil/backend/pkg/mq"
)

type GameCtrlUseCase struct {
	MessageService mq.MessageService
}

func (uc *GameCtrlUseCase) Pong(msg *mq.Message) {
	var value map[string]string
	json.Unmarshal(msg.Payload, &value)
	resp := map[string]string{
		"hi": "hi",
	}
	uc.MessageService.Response(*msg, resp)
}
