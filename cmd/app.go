package main

import (
	"github.com/joho/godotenv"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/oraksil/orakki/internal/presenter/di"
	"github.com/oraksil/orakki/internal/presenter/mq/handlers"
	"github.com/sangwonl/mqrpc"
)

func setupRoutes(mqSvc *mqrpc.MqService, routes []handlers.Route) {
	for _, r := range routes {
		mqSvc.AddHandler(r.MsgType, r.Handler)
	}
}

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	di.Resolve(func(
		serviceConf *services.ServiceConfig,
		mqSvc *mqrpc.MqService,
		setupHandler *handlers.SetupHandler,
		gamingHandler *handlers.GamingHandler) {

		setupRoutes(mqSvc, setupHandler.Routes())
		setupRoutes(mqSvc, gamingHandler.Routes())

		mqSvc.Run(false)
	})
}
