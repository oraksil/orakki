package usecases

import (
	"fmt"
	"time"

	"github.com/oraksil/orakki/internal/domain/engine"
)

type GamingUseCase struct {
	// ServiceConfig  *services.ServiceConfig
	// MessageService services.MessageService
	EngineFactory engine.EngineFactory

	gameEngine *engine.GameEngine
}

func (uc *GamingUseCase) StartGame() {
	go func() {
		for {
			fmt.Println("waiting game context is setup.")
			time.Sleep(1 * time.Second)
			if uc.EngineFactory.CanCreateEngine() {
				break
			}
		}

		fmt.Println("creating game engine.")
		uc.gameEngine = uc.EngineFactory.CreateEngine()

		fmt.Println("run game engine.")
		uc.gameEngine.Run()
	}()
}
