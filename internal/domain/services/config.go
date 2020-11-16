package services

type ServiceConfig struct {
	MqRpcUri        string
	MqRpcNamespace  string
	MqRpcIdentifier string

	GipanImageFramesIpcUri string
	GipanSoundFramesIpcUri string
	GipanCmdInputsIpcUri   string

	TurnServerUri      string
	TurnServerUsername string
	TurnServerPassword string
}
