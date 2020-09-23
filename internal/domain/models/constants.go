package models

const (
	MsgPrepareOrakki      = "MsgPrepareOrakki"
	MsgSetupWithNewOffer  = "MsgSetupWithNewOffer"
	MsgRemoteIceCandidate = "MsgRemoteIceCandidate"
	MsgStartGame          = "MsgStartGame"
)

const (
	OrakkiStateInit = iota
	OrakkiStateReady
	OrakkiStatePause
	OrakkiStatePlay
	OrakkiStateExit
	OrakkiStatePanic
)
