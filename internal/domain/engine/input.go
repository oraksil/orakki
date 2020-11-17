package engine

const (
	InputEventTypeSessionOpen  = 1
	InputEventTypeSessionClose = 2
	InputEventTypeHealthCheck  = 3
	InputEventTypeKeyMessage   = 4
)

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
			e.updatePlayerLiveness(in.PlayerId)

		case InputEventTypeKeyMessage:
			if slotNo, ok := e.playerSlots[in.PlayerId]; ok {
				e.gipan.WriteKeyInput(slotNo, in.Data)
			}
			e.updatePlayerLastInput()
		}

		if e.poisonPill {
			break
		}
	}
}
