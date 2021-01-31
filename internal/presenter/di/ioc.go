package di

import (
	"github.com/golobby/container"
)

func InitContainer() {
	container.Singleton(newServiceConfig)
	container.Singleton(newMqService)
	container.Singleton(newMessageService)
	container.Singleton(newWebRTCSession)
	container.Singleton(newGipanDriver)
	container.Singleton(newEngineFactory)
	container.Singleton(newSetupUseCase)
	container.Singleton(newGamingUseCase)
	container.Singleton(newSetupHandler)
	container.Singleton(newGamingHandler)
}

func Resolve(receiver interface{}) {
	container.Make(receiver)
}
