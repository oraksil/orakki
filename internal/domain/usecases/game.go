package usecases

import "gitlab.com/oraksil/orakki/internal/domain/services"

type GameCtrlUseCase struct {
	MessageService services.MessageService
}

func (uc *GameCtrlUseCase) Pong() interface{} {
	return map[string]string{
		"hi": "hi",
	}
}
