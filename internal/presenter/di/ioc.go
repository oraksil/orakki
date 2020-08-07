package di

import (
	"github.com/golobby/container"
	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/orakki/internal/domain/services"
	"gitlab.com/oraksil/orakki/internal/presenter/mq/handlers"
)

func InitContainer() {
	container.Singleton(newServiceConfig)
	container.Singleton(newMqService)
	container.Singleton(newMessageService)
	container.Singleton(newSystemMonitorUseCase)
	container.Singleton(newSystemHandler)
}

func InjectServiceConfig() *services.ServiceConfig {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)
	return serviceConf
}

func InjectMqService() *mqrpc.MqService {
	var svc *mqrpc.MqService
	container.Make(&svc)
	return svc
}

func InjectSystemHandler() *handlers.SystemHandler {
	var handler *handlers.SystemHandler
	container.Make(&handler)
	return handler
}
