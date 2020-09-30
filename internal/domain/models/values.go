package models

type PrepareOrakki struct {
	GameId int64
}

type GameInfo struct {
	GameId     int64
	MaxPlayers int
}

type PlayerParticipation struct {
	GameId   int64
	PlayerId int64
}

type Orakki struct {
	Id    string
	State int
}

type SdpInfo struct {
	SrcPeerId        int64 // game-id or player-id
	DstPeerId        int64 // player-id or game-id
	SdpBase64Encoded string
}

type IceCandidate struct {
	SrcPeerId        int64 // game-id or player-id
	DstPeerId        int64 // player-id or game-id
	IceBase64Encoded string
	Seq              int64
}
