package services

type ServiceConfig struct {
	MqRpcUri        string
	MqRpcNamespace  string
	MqRpcIdentifier string

	GipanImageFramesIpcUri string
	GipanSoundFramesIpcUri string
	GipanCmdInputsIpcUri   string

	TurnServerUri       string
	TurnServerSecretKey string
	TurnServerTTL       int

	PlayerHealthCheckTimeout  int
	PlayerHealthCheckInterval int
	PlayerIdleCheckTimeout    int
	PlayerIdleCheckInterval   int
}
