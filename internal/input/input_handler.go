package input

import (
	"log"

	"github.com/op/go-nanomsg"
)

type InputReader interface {
	ReadInput() ([]byte, error)
	AddInput(message []byte)
}

type InputHandler struct {
	Reader  InputReader
	Message chan []byte
	socket  *nanomsg.PushSocket
}

func (i *InputHandler) SetInputReader(reader InputReader) {
	i.Reader = reader
}

func CreateInputHandler(keyEvtQueuePath string) *InputHandler {
	var sock *nanomsg.PushSocket = nil

	if len(keyEvtQueuePath) > 0 {
		var err error = nil
		sock, err = nanomsg.NewPushSocket()
		if err != nil {
			panic(err)
		}

		_, err = sock.Connect(keyEvtQueuePath)
		if err != nil {
			panic(err)
		}
	}

	return &InputHandler{
		Message: make(chan []byte),
		Reader:  nil,
		socket:  sock,
	}
}

func (c *InputHandler) reader() {
	for {
		if c.Reader != nil {
			message, err := c.Reader.ReadInput()

			if err != nil {
				log.Println(err)
			} else {
				c.Message <- message
			}
		}
	}
}

func (c *InputHandler) writer_socket() {
	for {
		select {
		case message := <-c.Message:
			c.socket.Send(message, nanomsg.DontWait)
		}
	}
}

func (c *InputHandler) Run() {
	go c.reader()

	go c.writer_socket()
}

type WebRTCReader struct {
	messageChannel chan []byte
}

func CreateReader(bufferSize int) *WebRTCReader {
	r := WebRTCReader{
		make(chan []byte, bufferSize),
	}

	return &r
}

func (r *WebRTCReader) ReadInput() ([]byte, error) {
	return <-r.messageChannel, nil
}

func (r *WebRTCReader) AddInput(message []byte) {
	r.messageChannel <- message
}
