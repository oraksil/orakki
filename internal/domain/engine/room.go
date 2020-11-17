package engine

import (
	"time"

	"github.com/oraksil/orakki/internal/domain/models"
)

const initialPlayerSlotNo = 0

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
			if e.poisonPill {
				ticker.Stop()
			}
		}
	}
}

func (e *GameEngine) updatePlayerLiveness(playerId int64) {
	e.playerLastPings[playerId] = time.Now().Unix()
}

func (e *GameEngine) handleIdlePlayers() {
	const idleTimeout int64 = 60 // in seconds
	const checkInterval = 3 * time.Second

	handleIdleness := func() {
		now := time.Now().Unix()
		if now-e.playerLastInput > idleTimeout*3 {
			e.fsm.Event("suspend")
		} else if now-e.playerLastInput > idleTimeout {
			e.fsm.Event("idle")
		} else {
			e.fsm.Event("resume")
		}
	}

	ticker := time.NewTicker(checkInterval)
	for {
		select {
		case <-ticker.C:
			handleIdleness()
			if e.poisonPill {
				ticker.Stop()
			}
		}
	}
}

func (e *GameEngine) updatePlayerLastInput() {
	e.playerLastInput = time.Now().Unix()
}

func (e *GameEngine) notifyPlayerParticipation(msgType string, playerId int64) {
	e.messageService(msgType, &models.PlayerParticipation{
		GameId:   e.gameInfo.GameId,
		PlayerId: playerId,
	})
}
