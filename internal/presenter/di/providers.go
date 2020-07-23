package di

import (
	"github.com/golobby/container"
	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/orakki/internal/domain/services"
	"gitlab.com/oraksil/orakki/internal/domain/usecases"
	"gitlab.com/oraksil/orakki/internal/presenter/mq/handlers"
)

func newMqService() *mqrpc.MqService {
	return mqrpc.NewMqService("amqp://oraksil:oraksil@localhost:5672/", "oraksil")
}

func newMessageService() services.MessageService {
	var mqService *mqrpc.MqService
	container.Make(&mqService)

	return &mqrpc.DefaultMessageServiceImpl{MqService: mqService}
}

func newGameCtrlUseCase() *usecases.GameCtrlUseCase {
	var msgService mqrpc.MessageService
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
