package impl

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"encoding/json"

	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/pkg/utils"
)

const (
	videoFrameSizeWidth  = 480
	videoFrameSizeHeight = 320
	videoFps             = 23

	audioSampleSize     = 8 * 1024
	audioMaxSampleCount = 512
)

type KeyPad struct {
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

var baseKeyPad = KeyPad{
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

var gipanKeys = map[int]KeyPad{
	0: {
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
	1: {
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
	2: {
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
	3: {
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

type GipanDriverImpl struct {
	videoFrameBuffer utils.IpcBuffer
	audioFrameBuffer utils.IpcBuffer
	cmdInputBuffer   utils.IpcBuffer
}

func (g *GipanDriverImpl) ReadVideoFrame() ([]byte, error) {
	return g.videoFrameBuffer.Read()
}

func (g *GipanDriverImpl) ReadAudioFrame() ([]byte, error) {
	return g.audioFrameBuffer.Read()
}

type GipanCmd struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func (g *GipanDriverImpl) WriteKeyInput(playerSlotNo int, key []byte) error {
	// key input from browser: 005d, 005u

	keyCode, _ := strconv.Atoi(string(key[0:3]))
	keyState := key[3]

	basePad := reflect.ValueOf(baseKeyPad)
	gipanPad := reflect.ValueOf(gipanKeys[playerSlotNo])

	gipanKeyStr := ""
	for i := 0; i < basePad.NumField(); i++ {
		if basePad.Field(i).Int() == int64(keyCode) {
			keyName := basePad.Type().Field(i).Name
			gipanKeyCode := gipanPad.FieldByName(keyName).Int()
			gipanKeyStr = fmt.Sprintf("%03d%c", gipanKeyCode, keyState)
			break
		}
	}

	if gipanKeyStr == "" {
		return errors.New(fmt.Sprintf("failed to map key(%s) to gipan key.", string(key)))
	}

	return g.WriteCommand("key", []string{gipanKeyStr})
}

func (g *GipanDriverImpl) WriteCommand(cmd string, args []string) error {
	gipanCmdJsonStr, err := json.Marshal(GipanCmd{Cmd: cmd, Args: args})
	if err != nil {
		return err
	}

	return g.cmdInputBuffer.Write([]byte(gipanCmdJsonStr))
}

func NewGipanDriver(imagesIpcPath, soundsIpcPath, cmdsIpcPath string) engine.GipanDriver {
	maxVideoFrameBuffer := videoFrameSizeWidth * videoFrameSizeHeight * videoFps
	vb, err := utils.NewIpcBufferForRead(imagesIpcPath, maxVideoFrameBuffer)
	if err != nil {
		return nil
	}

	maxAudioFrameBuffer := audioSampleSize * audioMaxSampleCount
	ab, err := utils.NewIpcBufferForRead(soundsIpcPath, maxAudioFrameBuffer)
	if err != nil {
		return nil
	}

	cb, err := utils.NewIpcBufferForWrite(cmdsIpcPath)
	if err != nil {
		return nil
	}

	return &GipanDriverImpl{
		videoFrameBuffer: vb,
		audioFrameBuffer: ab,
		cmdInputBuffer:   cb,
	}
}
