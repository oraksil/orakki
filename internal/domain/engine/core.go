package engine

import (
	"time"

	"github.com/oraksil/orakki/internal/domain/models"
)

const (
	InputEventTypeSessionOpen  = 1
	InputEventTypeSessionClose = 2
	InputEventTypeHealthCheck  = 3
	InputEventTypeKeyMessage   = 4
)

type InputEvent struct {
	PlayerId int64
	Type     int
	Data     []byte
}

type RenderContext interface {
	WriteAudioFrame(buf []byte) error
	WriteVideoFrame(buf []byte) error
}

type InputContext interface {
	FetchInput() (InputEvent, error)
}

type SessionContext interface {
	CloseSession(playerId int64) error
}

type FrontInterface interface {
	WriteAudioFrame(buf []byte) error
	WriteVideoFrame(buf []byte) error
	FetchInput() (InputEvent, error)
	CloseSession(playerId int64) error
}

type GipanDriver interface {
	ReadAudioFrame() ([]byte, error)
	ReadVideoFrame() ([]byte, error)

	WriteKeyInput(playerSlotNo int, key []byte) error
}

const initialPlayerSlotNo = 0

type GameEngine struct {
	front          FrontInterface
	gipan          GipanDriver
	messageService func(msgType string, payload interface{})

	gameInfo        *models.GameInfo
	playerSlots     map[int64]int   // playerId: slotNo
	playerLastPings map[int64]int64 // playerId: unix time in secs

	running bool
}

func (e *GameEngine) Run(gameInfo *models.GameInfo, msgService func(string, interface{})) {
	// gipan -> renderer
	go e.handleAudioFrame()
	go e.handleVideoFrame()

	// input -> gipan
	go e.handleInputEvent()

	// kick unhealthy players
	go e.handleUnhealthyPlayers()

	e.messageService = msgService
	e.gameInfo = gameInfo
	e.running = true
}

func (e *GameEngine) handleAudioFrame() {
	for {
		buf, err := e.gipan.ReadAudioFrame()
		if err != nil {
			continue
		}

		e.front.WriteAudioFrame(buf)
	}
}

func (e *GameEngine) handleVideoFrame() {
	for {
		buf, err := e.gipan.ReadVideoFrame()
		if err != nil {
			continue
		}

		e.front.WriteVideoFrame(buf)
	}
}

func (e *GameEngine) handleInputEvent() {
	for {
		in, err := e.front.FetchInput()
		if err != nil {
			continue
		}

		switch inType := in.Type; inType {
		case InputEventTypeSessionOpen:
			e.joinPlayer(in.PlayerId)

		case InputEventTypeSessionClose:
			e.leavePlayer(in.PlayerId)

		case InputEventTypeHealthCheck:
			e.checkPlayerLiveness(in.PlayerId)

		case InputEventTypeKeyMessage:
			if slotNo, ok := e.playerSlots[in.PlayerId]; ok {
				e.gipan.WriteKeyInput(slotNo, in.Data)
			}
		}
	}
}

func (e *GameEngine) handleUnhealthyPlayers() {
	const unhealthyTimeout int64 = 20 // in seconds
	const checkInterval = 5 * time.Second

	kickUnhealthyPlayers := func() {
		now := time.Now().Unix()
		for playerId, lastPing := range e.playerLastPings {
			if now-lastPing > unhealthyTimeout {
				// kick player by closing channel
				e.front.CloseSession(playerId)
				e.leavePlayer(playerId)
			}
		}
	}

	ticker := time.NewTicker(checkInterval)
	for {
		select {
		case <-ticker.C:
			kickUnhealthyPlayers()
		}
	}
}

func isSlotOccupied(slotNumbers []int, slotNo int) bool {
	for _, i := range slotNumbers {
		if i == slotNo {
			return true
		}
	}
	return false
}

func (e *GameEngine) checkPlayerLiveness(playerId int64) {
	e.playerLastPings[playerId] = time.Now().Unix()
}

func (e *GameEngine) joinPlayer(playerId int64) {
	numOccupiedSlots := len(e.playerSlots)
	if numOccupiedSlots >= e.gameInfo.MaxPlayers {
		e.notifyPlayerParticipation(models.MsgPlayerJoinFailed, playerId)
		return
	}

	occupiedSlots := make([]int, 0, numOccupiedSlots)
	for _, slotNo := range e.playerSlots {
		occupiedSlots = append(occupiedSlots, slotNo)
	}

	for slotNo := initialPlayerSlotNo; slotNo < e.gameInfo.MaxPlayers; slotNo++ {
		if !isSlotOccupied(occupiedSlots, slotNo) {
			e.playerSlots[playerId] = slotNo
			e.playerLastPings[playerId] = time.Now().Unix()
			e.notifyPlayerParticipation(models.MsgPlayerJoined, playerId)
			return
		}
	}

	e.notifyPlayerParticipation(models.MsgPlayerJoinFailed, playerId)
}

func (e *GameEngine) leavePlayer(playerId int64) {
	if _, ok := e.playerSlots[playerId]; ok {
		delete(e.playerSlots, playerId)
		delete(e.playerLastPings, playerId)
		e.notifyPlayerParticipation(models.MsgPlayerLeft, playerId)
	}
}

func (e *GameEngine) notifyPlayerParticipation(msgType string, playerId int64) {
	e.messageService(msgType, &models.PlayerParticipation{
		GameId:   e.gameInfo.GameId,
		PlayerId: playerId,
	})
}

func NewGameEngine(f FrontInterface, g GipanDriver) *GameEngine {
	return &GameEngine{
		front:           f,
		gipan:           g,
		playerSlots:     make(map[int64]int),
		playerLastPings: make(map[int64]int64),
	}
}
