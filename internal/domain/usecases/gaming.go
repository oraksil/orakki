package usecases

import (
	"github.com/oraksil/orakki/internal/domain/engine"
)

type GamingUseCase struct {
	// ServiceConfig  *services.ServiceConfig
	// MessageService services.MessageService
	EngineFactory engine.EngineFactory

	gameEngine *engine.GameEngine
}

func (uc *GamingUseCase) StartGame(gameId string) {
	uc.gameEngine = uc.EngineFactory.CreateEngine()

	uc.gameEngine.Run()
}
