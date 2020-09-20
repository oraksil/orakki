package services

type ServiceConfig struct {
	MqRpcUri        string
	MqRpcNamespace  string
	MqRpcIdentifier string

	GipanImageFramesIpcPath string
	GipanSoundFramesIpcPath string
	GipanKeyInputsIpcPath   string
}
