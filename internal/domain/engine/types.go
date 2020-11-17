package engine

type EngineProps struct {
	PlayerHealthCheckTimeout  int
	PlayerHealthCheckInterval int
	PlayerIdleCheckTimeout    int
	PlayerIdleCheckInterval   int
}

type InputEvent struct {
	PlayerId int64
	Type     int
	Data     []byte
}

type RenderContext interface {
	WriteAudioFrame(buf []byte) error
	WriteVideoFrame(buf []byte) error
}

type InputContext interface {
	FetchInput() (InputEvent, error)
}

type SessionContext interface {
	CloseSession(playerId int64) error
}

type FrontInterface interface {
	WriteAudioFrame(buf []byte) error
	WriteVideoFrame(buf []byte) error
	FetchInput() (InputEvent, error)
	CloseSession(playerId int64) error
}

type GipanDriver interface {
	ReadAudioFrame() ([]byte, error)
	ReadVideoFrame() ([]byte, error)

	WriteKeyInput(playerSlotNo int, key []byte) error
	WriteCommand(cmd string, args []string) error
}
