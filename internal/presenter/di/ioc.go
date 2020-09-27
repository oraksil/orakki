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
	container.Singleton(newSetupUseCase)
	container.Singleton(newGamingUseCase)
	container.Singleton(newSetupHandler)
	container.Singleton(newGamingHandler)
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

func InjectSetupHandler() *handlers.SetupHandler {
	var handler *handlers.SetupHandler
	container.Make(&handler)
	return handler
}

func InjectGamingHandler() *handlers.GamingHandler {
	var handler *handlers.GamingHandler
	container.Make(&handler)
	return handler
}
