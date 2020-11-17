package engine

import "fmt"

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

func (e *GameEngine) pauseGipan() {
	fmt.Println("pausing gipan")
	e.gipan.WriteCommand("ctrl", []string{"pause"})
}

func (e *GameEngine) resumeGipan() {
	fmt.Println("resuming gipan")
	e.gipan.WriteCommand("ctrl", []string{"resume"})
}
