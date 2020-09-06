package services

import "time"

type MessageService interface {
	Identifier() string
	Send(to, msgType string, payload interface{}) error
	SendToAny(msgType string, payload interface{}) error
	Broadcast(msgType string, payload interface{}) error
	Request(to, msgType string, payload interface{}, timeout time.Duration) (interface{}, error)
}
