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

	gameInfo        *models.GameInfo
	playerSlots     map[int64]int   // playerId: slotNo
	playerLastPings map[int64]int64 // playerId: unix time in secs
	playerLastInput int64

	poisonPill bool

	fsm *fsm.FSM
}

func (e *GameEngine) Reset() {
	e.poisonPill = true
}

func (e *GameEngine) Run(gameInfo *models.GameInfo, msgService func(string, interface{})) {
	// gipan -> renderer
	go e.handleAudioFrame()
	go e.handleVideoFrame()

	// input -> gipan
	go e.handleInputEvent()

	// handle unhealthy or idle players
	go e.handleUnhealthyPlayers()
	go e.handleIdlePlayers()
	e.updatePlayerLastInput()

	e.messageService = msgService
	e.gameInfo = gameInfo
	e.poisonPill = false

	e.fsm = e.newFSM()
	e.fsm.Event("start")
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
