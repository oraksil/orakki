package di

import (
	"github.com/golobby/container"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/oraksil/orakki/internal/presenter/mq/handlers"
	"github.com/sangwonl/mqrpc"
)

func InitContainer() {
	container.Singleton(newServiceConfig)
	container.Singleton(newMqService)
	container.Singleton(newMessageService)
	container.Singleton(newWebRTCSession)
	container.Singleton(newEngineFactory)
	container.Singleton(newSystemUseCase)
	container.Singleton(newSetupUseCase)
	container.Singleton(newSystemHandler)
	container.Singleton(newSetupHandler)
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

func InjectSetupHandler() *handlers.SetupHandler {
	var handler *handlers.SetupHandler
	container.Make(&handler)
	return handler
}
