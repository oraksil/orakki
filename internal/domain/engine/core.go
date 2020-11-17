package engine

import (
	"fmt"
	"os"

	"github.com/looplab/fsm"
	"github.com/oraksil/orakki/internal/domain/models"
)

type GameEngine struct {
	front          FrontInterface
	gipan          GipanDriver
	messageService func(msgType string, payload interface{})

	// props
	props    *EngineProps
	gameInfo *models.GameInfo

	// state
	fsm        *fsm.FSM
	poisonPill bool

	// players
	playerSlots     map[int64]int   // playerId: slotNo
	playerLastPings map[int64]int64 // playerId: unix time in secs
	playerLastInput int64
}

func (e *GameEngine) Reset() {
	e.poisonPill = true
}

func (e *GameEngine) Run(props *EngineProps, gameInfo *models.GameInfo, msgService func(string, interface{})) {
	e.messageService = msgService

	e.props = props
	e.gameInfo = gameInfo

	e.fsm = e.newFSM()
	e.fsm.Event("start")
	e.poisonPill = false

	e.updatePlayerLastInput()

	// gipan -> renderer
	go e.handleAudioFrame()
	go e.handleVideoFrame()

	// input -> gipan
	go e.handleInputEvent()

	// handle unhealthy or idle players
	go e.handleUnhealthyPlayers()
	go e.handleIdlePlayers()
}

func (e *GameEngine) shutdown() {
	fmt.Println("shutting down..")
	os.Exit(0)
}

func NewGameEngine(f FrontInterface, g GipanDriver) *GameEngine {
	return &GameEngine{
		front:           f,
		gipan:           g,
		playerSlots:     make(map[int64]int),
		playerLastPings: make(map[int64]int64),
	}
}
