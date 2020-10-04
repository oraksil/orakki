package models

const (
	MsgPrepareOrakki      = "MsgPrepareOrakki"
	MsgSetupWithNewOffer  = "MsgSetupWithNewOffer"
	MsgRemoteIceCandidate = "MsgRemoteIceCandidate"
	MsgStartGame          = "MsgStartGame"
	MsgPlayerJoined       = "MsgPlayerJoined"
	MsgPlayerJoinFailed   = "MsgPlayerJoinFailed"
	MsgPlayerLeft         = "MsgPlayerLeft"
)

const (
	OrakkiStateInit = iota
	OrakkiStateReady
	OrakkiStatePause
	OrakkiStatePlay
	OrakkiStateExit
	OrakkiStatePanic
)
