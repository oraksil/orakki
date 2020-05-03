package utils

import (
	"github.com/op/go-nanomsg"
)

type FrameBuffer interface {
	Open(path string, frameSize, maxBufSize int64) error
	Close()
	GetFrameSize() int64
	GetBuffer() []byte
}

type IpcPipelineBuffer struct {
	socket     *nanomsg.PullSocket
	frameSize  int64
	maxBufSize int64
}

func (b *IpcPipelineBuffer) Open(path string, frameSize, maxBufSize int64) error {
	sock, err := nanomsg.NewPullSocket()
	if err != nil {
		panic(err)
	}

	sock.SetRecvMaxSize(int64(maxBufSize))
	_, err = sock.Connect(path)
	if err != nil {
		panic(err)
	}

	b.socket = sock
	b.frameSize = frameSize
	b.maxBufSize = maxBufSize

	return nil
}

func (b *IpcPipelineBuffer) Close() {
	b.socket.Close()
}

func (b *IpcPipelineBuffer) GetFrameSize() int64 {
	return b.frameSize
}

func (b *IpcPipelineBuffer) GetBuffer() []byte {
	buf, err := b.socket.Recv(0)
	if err != nil {
		return nil
	}
	return buf
}
