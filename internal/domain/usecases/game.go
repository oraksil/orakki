package usecases

import (
	"gitlab.com/oraksil/sil/backend/pkg/mq"
)

type GameCtrlUseCase struct {
	MessageService mq.MessageService
}

func (uc *GameCtrlUseCase) Pong() interface{} {
	return map[string]string{
		"hi": "hi",
	}
}
