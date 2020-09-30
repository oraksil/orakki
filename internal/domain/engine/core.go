package engine

import (
	"github.com/oraksil/orakki/internal/domain/models"
)

type RenderContext interface {
	WriteAudioFrame(buf []byte) error
	WriteVideoFrame(buf []byte) error
}

type Renderer interface {
	WriteAudioFrame(buf []byte) error
	WriteVideoFrame(buf []byte) error
}

const (
	InputEventTypeSessionOpen  = 1
	InputEventTypeSessionClose = 2
	InputEventTypeKeyMessage   = 3
)

type InputEvent struct {
	PlayerId int64
	Type     int
	Data     []byte
}

type InputContext interface {
	FetchInput() (InputEvent, error)
}

type InputHandler interface {
	FetchInput() (InputEvent, error)
}

type GipanDriver interface {
	ReadAudioFrame() ([]byte, error)
	ReadVideoFrame() ([]byte, error)

	WriteKeyInput(playerSlotNo int, key []byte) error
}

const initialPlayerSlotNo = 0

type GameEngine struct {
	renderer       Renderer
	input          InputHandler
	gipan          GipanDriver
	messageService func(msgType string, payload interface{})

	gameInfo    *models.GameInfo
	playerSlots map[int64]int

	running bool
}

func (e *GameEngine) Run(gameInfo *models.GameInfo, msgService func(string, interface{})) {
	// gipan -> renderer
	go e.handleAudioFrame()
	go e.handleVideoFrame()

	// input -> gipan
	go e.handleInputEvent()

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

		e.renderer.WriteAudioFrame(buf)
	}
}

func (e *GameEngine) handleVideoFrame() {
	for {
		buf, err := e.gipan.ReadVideoFrame()
		if err != nil {
			continue
		}

		e.renderer.WriteVideoFrame(buf)
	}
}

func (e *GameEngine) handleInputEvent() {
	for {
		in, err := e.input.FetchInput()
		if err != nil {
			continue
		}

		switch inType := in.Type; inType {
		case InputEventTypeSessionOpen:
			e.joinPlayer(in.PlayerId)

		case InputEventTypeSessionClose:
			e.leavePlayer(in.PlayerId)

		case InputEventTypeKeyMessage:
			if slotNo, ok := e.playerSlots[in.PlayerId]; ok {
				e.gipan.WriteKeyInput(slotNo, in.Data)
			}
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

func (e *GameEngine) joinPlayer(playerId int64) {
	numOccupiedSlots := len(e.playerSlots)
	if numOccupiedSlots >= e.gameInfo.MaxPlayers {
		e.notifyPlayerJoinFailed(playerId)
		return
	}

	occupiedSlots := make([]int, 0, numOccupiedSlots)
	for _, slotNo := range e.playerSlots {
		occupiedSlots = append(occupiedSlots, slotNo)
	}

	for slotNo := initialPlayerSlotNo; slotNo < e.gameInfo.MaxPlayers; slotNo++ {
		if !isSlotOccupied(occupiedSlots, slotNo) {
			e.playerSlots[playerId] = slotNo
			e.notifyPlayerJoined(playerId)
			return
		}
	}

	e.notifyPlayerJoinFailed(playerId)
}

func (e *GameEngine) leavePlayer(playerId int64) {
	if _, ok := e.playerSlots[playerId]; ok {
		delete(e.playerSlots, playerId)
		e.notifyPlayerLeft(playerId)
	}
}

func (e *GameEngine) notifyPlayerJoined(playerId int64) {
	e.messageService(models.MsgPlayerJoined, &models.PlayerParticipation{
		GameId:   e.gameInfo.GameId,
		PlayerId: playerId,
	})
}

func (e *GameEngine) notifyPlayerJoinFailed(playerId int64) {
	e.messageService(models.MsgPlayerJoinFailed, &models.PlayerParticipation{
		GameId:   e.gameInfo.GameId,
		PlayerId: playerId,
	})
}

func (e *GameEngine) notifyPlayerLeft(playerId int64) {
	e.messageService(models.MsgPlayerLeft, &models.PlayerParticipation{
		GameId:   e.gameInfo.GameId,
		PlayerId: playerId,
	})
}

func NewGameEngine(r Renderer, i InputHandler, g GipanDriver) *GameEngine {
	return &GameEngine{
		renderer:    r,
		input:       i,
		gipan:       g,
		playerSlots: make(map[int64]int),
	}
}
