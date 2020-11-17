package usecases

import (
	"fmt"
	"time"

	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
)

type GamingUseCase struct {
	ServiceConfig  *services.ServiceConfig
	MessageService services.MessageService
	EngineFactory  engine.EngineFactory

	gameEngine *engine.GameEngine
}

func (uc *GamingUseCase) StartGame(gameInfo *models.GameInfo) {
	engineProps := engine.EngineProps{
		PlayerHealthCheckTimeout:  uc.ServiceConfig.PlayerHealthCheckTimeout,
		PlayerHealthCheckInterval: uc.ServiceConfig.PlayerHealthCheckInterval,
		PlayerIdleCheckTimeout:    uc.ServiceConfig.PlayerIdleCheckTimeout,
		PlayerIdleCheckInterval:   uc.ServiceConfig.PlayerIdleCheckInterval,
	}

	msgService := func(msgType string, payload interface{}) {
		uc.MessageService.SendToAny(msgType, payload)
	}

	go func() {
		for {
			fmt.Println("waiting game context is setup.")
			time.Sleep(1 * time.Second)
			if uc.EngineFactory.CanCreateEngine() {
				break
			}
		}

		if uc.gameEngine != nil {
			fmt.Println("resetting game engine.")
			uc.gameEngine.Reset()
		}

		fmt.Println("creating game engine.")
		uc.gameEngine = uc.EngineFactory.CreateEngine()

		fmt.Println("run game engine.")
		uc.gameEngine.Run(&engineProps, gameInfo, msgService)
	}()
}
