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

type PeerInfo struct {
	Token    string
	GameId   int64
	PlayerId int64
}

type SdpInfo struct {
	Peer             PeerInfo
	SdpBase64Encoded string
}

type IceCandidate struct {
	Peer             PeerInfo
	IceBase64Encoded string
	Seq              int64
}
