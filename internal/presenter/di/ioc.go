package di

import (
	"github.com/golobby/container"
	"gitlab.com/oraksil/orakki/internal/presenter/mq/handlers"
	"gitlab.com/oraksil/sil/backend/pkg/mq"
)

func InitContainer() {
	container.Singleton(newMqService)
	container.Singleton(newMessageService)
	container.Singleton(newGameCtrlUseCase)
	container.Singleton(newHelloHandler)
}

func InjectMqService() *mq.MqService {
	var svc *mq.MqService
	container.Make(&svc)
	return svc
}

func InjectHelloHandler() *handlers.HelloHandler {
	var handler *handlers.HelloHandler
	container.Make(&handler)
	return handler
}
