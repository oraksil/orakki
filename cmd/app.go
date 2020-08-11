package main

import (
	"github.com/joho/godotenv"
	"gitlab.com/oraksil/orakki/internal/presenter/di"
)

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectSystemHandler())

	conf := di.InjectServiceConfig()
	mqSvc.Run(conf.PeerName)
}
