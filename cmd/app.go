package main

import (
	"gitlab.com/oraksil/orakki/internal/presenter/di"
)

func main() {
	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectHelloHandler())
	mqSvc.Run("orakki-temp")
}
