package engine

func (e *GameEngine) handleAudioFrame() {
	for {
		buf, err := e.gipan.ReadAudioFrame()
		if err != nil {
			continue
		}

		e.front.WriteAudioFrame(buf)

		if e.poisonPill {
			break
		}
	}
}

func (e *GameEngine) handleVideoFrame() {
	for {
		buf, err := e.gipan.ReadVideoFrame()
		if err != nil {
			continue
		}

		e.front.WriteVideoFrame(buf)

		if e.poisonPill {
			break
		}
	}
}

func (e *GameEngine) handleKeyInput(in InputEvent) {
	if slotNo, ok := e.playerSlots[in.PlayerId]; ok {
		e.gipan.WriteKeyInput(slotNo, in.Data)
	}
}
