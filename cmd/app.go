package main

import (
	"github.com/joho/godotenv"
	"github.com/oraksil/orakki/internal/presenter/di"
)

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectSystemHandler())
	mqSvc.AddHandler(di.InjectSetupHandler())

	conf := di.InjectServiceConfig()
	mqSvc.Run(conf.PeerName)
}
