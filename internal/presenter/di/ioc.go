package di

import (
	"github.com/golobby/container"
	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/orakki/internal/presenter/mq/handlers"
)

func InitContainer() {
	container.Singleton(newMqService)
	container.Singleton(newMessageService)
	container.Singleton(newGameCtrlUseCase)
	container.Singleton(newHelloHandler)
}

func InjectMqService() *mqrpc.MqService {
	var svc *mqrpc.MqService
	container.Make(&svc)
	return svc
}

func InjectHelloHandler() *handlers.HelloHandler {
	var handler *handlers.HelloHandler
	container.Make(&handler)
	return handler
}
