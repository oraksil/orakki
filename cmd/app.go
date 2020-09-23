package main

import (
	"github.com/joho/godotenv"
	"github.com/oraksil/orakki/internal/presenter/di"
)

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectSetupHandler())
	mqSvc.AddHandler(di.InjectGamingHandler())

	conf := di.InjectServiceConfig()
	mqSvc.Run(conf.MqRpcIdentifier, false)
}
