package di

import (
	"github.com/golobby/container"
	"gitlab.com/oraksil/orakki/internal/domain/usecases"
	"gitlab.com/oraksil/orakki/internal/presenter/mq/handlers"
	"gitlab.com/oraksil/sil/backend/pkg/mq"
)

func newMqService() *mq.MqService {
	return mq.NewMqService("amqp://oraksil:oraksil@localhost:5672/", "oraksil.mq.p2p", "oraksil.mq.broadcast")
}

func newMessageService() mq.MessageService {
	var mqService *mq.MqService
	container.Make(&mqService)

	return &mq.DefaultMessageServiceImpl{MqService: mqService}
}

func newGameCtrlUseCase() *usecases.GameCtrlUseCase {
	var msgService mq.MessageService
	container.Make(&msgService)

	return &usecases.GameCtrlUseCase{MessageService: msgService}
}

func newHelloHandler() *handlers.HelloHandler {
	var gameCtrlUseCase *usecases.GameCtrlUseCase
	container.Make(&gameCtrlUseCase)

	return &handlers.HelloHandler{
		GameCtrlUseCase: gameCtrlUseCase,
	}
}
