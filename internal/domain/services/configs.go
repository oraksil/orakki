package services

type ServiceConfig struct {
	MqRpcUri       string
	MqRpcNamespace string

	UseStaticOrakki bool

	OrakkiId string
	PeerName string
}
