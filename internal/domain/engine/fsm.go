package engine

import (
	"github.com/looplab/fsm"
)

func (e *GameEngine) newFSM() *fsm.FSM {
	// *ready
	//    |          +------(idle)-------+
	// (start)       |                   |
	//    |          |                   v
	//    +-----> running --(pause)--> paused ---+
	//               ^                   |       |
	//               |                   |       |
	//               +------(resume)-----+       |
	//                                           |
	//                                           |
	// *shutdown <----------(suspend)------------+

	return fsm.NewFSM("ready",
		fsm.Events{
			{Name: "start", Src: []string{"ready"}, Dst: "running"},
			{Name: "pause", Src: []string{"running"}, Dst: "paused"},
			{Name: "idle", Src: []string{"running"}, Dst: "paused"},
			{Name: "resume", Src: []string{"paused"}, Dst: "running"},
			{Name: "suspend", Src: []string{"paused"}, Dst: "shutdown"},
		},
		fsm.Callbacks{
			"paused":   func(ev *fsm.Event) { e.pauseGipan() },
			"running":  func(ev *fsm.Event) { e.resumeGipan() },
			"shutdown": func(ev *fsm.Event) { e.shutdown() },
		})
}
