package models

type PrepareOrakki struct {
	GameId int64
}

type Orakki struct {
	Id    string
	State int
}

type SdpInfo struct {
	PeerId           int64 // game-id or player-id
	SdpBase64Encoded string
}

type IceCandidate struct {
	PeerId           int64 // game-id or player-id
	IceBase64Encoded string
}
