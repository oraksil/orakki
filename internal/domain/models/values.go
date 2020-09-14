package models

type OrakkiState struct {
	OrakkiId string
	State    int
}

type SdpInfo struct {
	PeerId           string
	SdpBase64Encoded string
}

type IceCandidate struct {
	PeerId           string
	IceBase64Encoded string
}

type SetupAnswer struct {
	PeerId string
	Answer string
}

type Icecandidate struct {
	PlayerId  int64
	OrakkiId  string
	IceString string
}
