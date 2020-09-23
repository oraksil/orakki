package models

const (
	MsgPrepareOrakki      = "MsgPrepareOrakki"
	MsgSetupWithNewOffer  = "MsgSetupWithNewOffer"
	MsgRemoteIceCandidate = "MsgRemoteIceCandidate"
)

const (
	OrakkiStateInit = iota
	OrakkiStateReady
	OrakkiStatePause
	OrakkiStatePlay
	OrakkiStateExit
	OrakkiStatePanic
)
