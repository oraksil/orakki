package main

import (
	"gitlab.com/oraksil/orakki/internal/presenter/di"
)

func main() {
	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectSystemHandler())

	conf := di.InjectServiceConfig()
	mqSvc.Run(conf.PeerName)
}
