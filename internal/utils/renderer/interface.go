package renderer

import (
	"gitlab.com/oraksil/orakki/internal/utils"
	"gitlab.com/oraksil/orakki/internal/utils/input"
)

type RendererType int

const (
	TypeWebRTCStreamRenderer RendererType = 1
)

type Renderer interface {
	StartWithFrameBuffer(fb utils.FrameBuffer, sb utils.FrameBuffer, ih *input.InputHandler)
}

type FrameMetaDTO struct {
	Fps int
	W   int
	H   int
}

type FrameInfo struct {
	w          int
	h          int
	fps        int
	colorDepth int
}

func CreateFrameInfo(w, h, fps int) FrameInfo {
	return FrameInfo{
		w:          w,
		h:          h,
		fps:        fps,
		colorDepth: 4,
	}
}

func (f FrameInfo) MaxFrameBufferSize() int64 {
	frame_cap := 10
	return int64(f.w * f.h * f.colorDepth * frame_cap)
}

func (f FrameInfo) SingleFrameSize() int64 {
	return int64(f.w * f.h * f.colorDepth)
}

func CreateRenderer(rendererType RendererType, frameInfo FrameInfo) Renderer {
	switch rendererType {
	case TypeWebRTCStreamRenderer:
		return createWebRTCStreamRenderer(frameInfo)
	}
	panic("not supported renderer type.")
}
