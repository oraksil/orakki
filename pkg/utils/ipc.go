package utils

import (
	"github.com/op/go-nanomsg"
)

type IpcBuffer interface {
	Read() ([]byte, error)
	Write(buf []byte) error
	Close()
}

type IpcReadBuffer struct {
	socket *nanomsg.PullSocket
}

func (b *IpcReadBuffer) Read() ([]byte, error) {
	buf, err := b.socket.Recv(0)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (b *IpcReadBuffer) Write(buf []byte) error {
	panic("write is not allowed for read buffer")
}

func (b *IpcReadBuffer) Close() {
	b.socket.Close()
}

type IpcWriteBuffer struct {
	socket *nanomsg.PushSocket
}

func (b *IpcWriteBuffer) Read() ([]byte, error) {
	panic("read is not allowed for write buffer")
}

func (b *IpcWriteBuffer) Write(buf []byte) error {
	_, err := b.socket.Send(buf, nanomsg.DontWait)
	if err != nil {
		return err
	}
	return nil
}

func (b *IpcWriteBuffer) Close() {
	b.socket.Close()
}

func NewIpcBufferForRead(path string, maxBufSize int) (IpcBuffer, error) {
	sock, err := nanomsg.NewPullSocket()
	if err != nil {
		return nil, err
	}

	sock.SetRecvMaxSize(int64(maxBufSize))
	_, err = sock.Connect(path)
	if err != nil {
		return nil, err
	}

	ipcBuf := &IpcReadBuffer{
		socket: sock,
	}

	return ipcBuf, nil
}

func NewIpcBufferForWrite(path string) (IpcBuffer, error) {
	sock, err := nanomsg.NewPushSocket()
	if err != nil {
		return nil, err
	}

	_, err = sock.Connect(path)
	if err != nil {
		return nil, err
	}

	ipcBuf := &IpcWriteBuffer{
		socket: sock,
	}

	return ipcBuf, nil
}
