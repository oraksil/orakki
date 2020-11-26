package usecases

import (
	"fmt"
	"os"
	"time"

	"github.com/looplab/fsm"
	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
)

type GamingUseCase struct {
	ServiceConfig  *services.ServiceConfig
	MessageService services.MessageService
	EngineFactory  engine.EngineFactory
	GipanDriver    engine.GipanDriver

	fsm    *fsm.FSM
	engine *engine.GameEngine
}

func (uc *GamingUseCase) StartGame(gameInfo *models.GameInfo) {
	uc.fsm = uc.newFSM()

	engineProps := engine.EngineProps{
		PlayerHealthCheckTimeout:  uc.ServiceConfig.PlayerHealthCheckTimeout,
		PlayerHealthCheckInterval: uc.ServiceConfig.PlayerHealthCheckInterval,
		PlayerIdleCheckTimeout:    uc.ServiceConfig.PlayerIdleCheckTimeout,
		PlayerIdleCheckInterval:   uc.ServiceConfig.PlayerIdleCheckInterval,
	}

	msgService := func(msgType string, payload interface{}) {
		uc.MessageService.SendToAny(msgType, payload)
	}

	engineEvtHandler := func(event string) {
		uc.fsm.Event(event)
	}

	idleTimeout := time.Duration(engineProps.PlayerIdleCheckTimeout) * time.Second
	time.AfterFunc(idleTimeout*3, func() {
		if uc.fsm.Current() == "ready" {
			uc.fsm.Event("suspend")
		}
	})

	go func() {
		for {
			fmt.Println("waiting game context is setup.")
			time.Sleep(1 * time.Second)
			if uc.EngineFactory.CanCreateEngine() {
				break
			}
		}

		if uc.engine != nil {
			fmt.Println("resetting game engine.")
			uc.engine.Reset()
		}

		fmt.Println("creating game engine.")
		uc.engine = uc.EngineFactory.CreateEngine()

		fmt.Println("run game engine.")
		uc.engine.Run(&engineProps, gameInfo, msgService, engineEvtHandler)

		uc.fsm.Event("start")
	}()
}

func (uc *GamingUseCase) newFSM() *fsm.FSM {
	//  ready---------------(suspend)------------------+
	//    |                                            |
	//    |                                            |
	//    |          +------(idle)-------+             |
	// (start)       |                   |             |
	//    |          |                   v             |
	//    +-----> running --(pause)--> paused ---+     |
	//               ^                   |       |     |
	//               |                   |       |     |
	//               +------(resume)-----+       |     |
	//                                           |     |
	//                                           |     |
	// *shutdown <----------(suspend)------------+ <---+

	return fsm.NewFSM("ready",
		fsm.Events{
			{Name: "suspend", Src: []string{"ready"}, Dst: "shutdown"},
			{Name: "start", Src: []string{"ready"}, Dst: "running"},
			{Name: "pause", Src: []string{"running"}, Dst: "paused"},
			{Name: "idle", Src: []string{"running"}, Dst: "paused"},
			{Name: "resume", Src: []string{"paused"}, Dst: "running"},
			{Name: "suspend", Src: []string{"paused"}, Dst: "shutdown"},
		},
		fsm.Callbacks{
			"paused":   func(ev *fsm.Event) { uc.pause() },
			"running":  func(ev *fsm.Event) { uc.resume() },
			"shutdown": func(ev *fsm.Event) { uc.shutdown() },
		})
}

func (uc *GamingUseCase) pause() {
	fmt.Println("pausing gipan")
	uc.GipanDriver.WriteCommand("ctrl", []string{"pause"})
}

func (uc *GamingUseCase) resume() {
	fmt.Println("resuming gipan")
	uc.GipanDriver.WriteCommand("ctrl", []string{"resume"})
}

func (uc *GamingUseCase) shutdown() {
	uc.GipanDriver.WriteCommand("ctrl", []string{"shutdown"})
	os.Exit(0)
}
