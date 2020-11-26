package engine

import (
	"github.com/oraksil/orakki/internal/domain/models"
)

type GameEngine struct {
	front          FrontInterface
	gipan          GipanDriver
	messageService func(msgType string, payload interface{})
	eventHandler   func(event string)

	// props
	props    *EngineProps
	gameInfo *models.GameInfo

	// state
	poisonPill bool

	// players
	playerSlots     map[int64]int   // playerId: slotNo
	playerLastPings map[int64]int64 // playerId: unix time in secs
	playerLastInput int64
}

func (e *GameEngine) Reset() {
	e.poisonPill = true
}

func (e *GameEngine) Run(
	props *EngineProps,
	gameInfo *models.GameInfo,
	msgService func(string, interface{}),
	eventHandler func(string)) {

	e.messageService = msgService
	e.eventHandler = eventHandler

	e.props = props
	e.gameInfo = gameInfo
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

func NewGameEngine(f FrontInterface, g GipanDriver) *GameEngine {
	return &GameEngine{
		front:           f,
		gipan:           g,
		playerSlots:     make(map[int64]int),
		playerLastPings: make(map[int64]int64),
	}
}
