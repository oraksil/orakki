package services

type MessageService interface {
	Identifier() string
	AllParticipants() []string

	Send(to string, msgType string, payload interface{})
	Broadcast(msgType string, payload interface{})
	Request(to string, msgType string, payload interface{}) interface{}
}
