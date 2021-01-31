package main

import (
	"github.com/joho/godotenv"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/oraksil/orakki/internal/presenter/di"
	"github.com/oraksil/orakki/internal/presenter/mq/handlers"
	"github.com/sangwonl/mqrpc"
)

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	di.Resolve(func(
		serviceConf *services.ServiceConfig,
		mqSvc *mqrpc.MqService,
		setupHandler *handlers.SetupHandler,
		gamingHandler *handlers.GamingHandler) {

		mqSvc.AddHandler(setupHandler)
		mqSvc.AddHandler(gamingHandler)

		mqSvc.Run(serviceConf.MqRpcIdentifier, false)
	})
}
