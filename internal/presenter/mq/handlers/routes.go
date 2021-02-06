package handlers

import (
	"github.com/sangwonl/mqrpc"
)

type Route struct {
	MsgType mqrpc.MsgType
	Handler mqrpc.HandlerFunc
}
