package input

import (
	"log"
	"fmt"
	"reflect"
	"strconv"

	"github.com/op/go-nanomsg"
)

const MAX_USERS int = 4

type CommandKey struct {
	up      int8
	down    int8
	left    int8
	right   int8
	coin    int8
	start   int8
	button1 int8
	button2 int8
	button3 int8
	button4 int8
	button5 int8
	button6 int8
	button7 int8
	button8 int8
}

var taktakiKey = CommandKey{
	up:      38,
	down:    40,
	left:    37,
	right:   39,
	coin:    49,
	start:   50,
	button1: 65,
	button2: 83,
	button3: 68,
	button4: 90,
	button5: 88,
	button6: 67,
	button7: 81,
	button8: 87,
}

var gipanKeys = map[int]CommandKey{
	0: CommandKey{
		up:      38,
		down:    40,
		left:    37,
		right:   39,
		coin:    53,
		start:   49,
		button1: 90,
		button2: 88,
		button3: 67,
		button4: 65,
		button5: 83,
		button6: 68,
		button7: 81,
		button8: 87,
	},
	1: CommandKey{
		up:      1,
		down:    2,
		left:    3,
		right:   4,
		coin:    6,
		start:   5,
		button1: 7,
		button2: 8,
		button3: 9,
		button4: 10,
		button5: 11,
		button6: 12,
		button7: 13,
		button8: 14,
	},
	2: CommandKey{
		up:      15,
		down:    16,
		left:    17,
		right:   18,
		coin:    20,
		start:   19,
		button1: 21,
		button2: 22,
		button3: 23,
		button4: 24,
		button5: 25,
		button6: 26,
		button7: 27,
		button8: 28,
	},
	3: CommandKey{
		up:      29,
		down:    30,
		left:    31,
		right:   32,
		coin:    43,
		start:   44,
		button1: 36,
		button2: 37,
		button3: 38,
		button4: 33,
		button5: 34,
		button6: 35,
		button7: 39,
		button8: 41,
	},
}

type InputReader interface {
	ReadInput() ([]byte, error)
	AddInput(message []byte)
}

type InputHandler struct {
	Reader  InputReader
	Message chan []byte
	socket  *nanomsg.PushSocket

		// gamepads
	padIDs map[int]*uint16
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
		padIDs:  make(map[int]*uint16),
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

func (c *InputHandler) assignPad(seq int, ID *uint16) {
	fmt.Printf("assign pad of %d to # %d\n", ID, seq)
	c.padIDs[seq] = ID
}

func (c *InputHandler) RemovePad(ID *uint16) bool {
	for i := 0; i < MAX_USERS; i++ {
		if c.padIDs[i] == ID {
			fmt.Printf("remove pad of %d from # %d\n", ID, i)
			c.padIDs[i] = nil
			return true
		}
	}

	return false
}

func (c *InputHandler) SetPadInOrder(ID *uint16) bool {
	for i := 0; i < MAX_USERS; i++ {
		if c.padIDs[i] == nil {
			c.assignPad(i, ID)
			return true
		}
	}

	return false
}

func (c *InputHandler) GetPlayerOrderFromPad(ID *uint16) int {
	for i := 0; i < MAX_USERS; i++ {
		if c.padIDs[i] == ID {
			return i
		}
	}

	return -1
}

func ConvertToGipanKeys(padID int, userKey []byte) []byte {
	userKeyCode, _ := strconv.Atoi(string(userKey[1:3]))

	pad0 := reflect.ValueOf(taktakiKey)
	userpad := reflect.ValueOf(gipanKeys[padID])

	for i := 0; i < pad0.NumField(); i++ {
		if pad0.Field(i).Int() == int64(userKeyCode) {
			commandType := pad0.Type().Field(i).Name
			ret := []byte(fmt.Sprintf("%03d%c", userpad.FieldByName(commandType).Int(), userKey[3]))

			return ret
		}
	}
	return nil
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
