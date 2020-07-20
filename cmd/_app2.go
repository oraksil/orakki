package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/oraksil/orakki/internal/input"
	"github.com/oraksil/orakki/internal/renderer"
	"github.com/oraksil/orakki/pkg/utils"
	"github.com/spf13/cobra"
)

var frameSource string
var soundSource string
var keyEvtQueue string
var resolution string
var fps string

func cmdRunApp() *cobra.Command {
	cmdShow := cobra.Command{
		Use:   "app [options]",
		Short: "Run main app",
		Args:  cobra.MinimumNArgs(0),
		Run:   run,
	}

	cmdShow.PersistentFlags().StringVar(&frameSource, "framesrc", "", "framesrc")
	cmdShow.PersistentFlags().StringVar(&soundSource, "soundsrc", "", "soundsrc")
	cmdShow.PersistentFlags().StringVar(&keyEvtQueue, "keyevtqueue", "", "keyevtqueue")
	cmdShow.PersistentFlags().StringVar(&resolution, "resolution", "480x320", "resolution")
	cmdShow.PersistentFlags().StringVar(&fps, "fps", "30", "fps")

	return &cmdShow
}

func createFrameBuffer(srcPath string, fi renderer.FrameInfo) utils.FrameBuffer {
	fb := &utils.IpcPipelineBuffer{}
	err := fb.Open(srcPath, fi.SingleFrameSize(), fi.MaxFrameBufferSize())
	if err != nil {
		panic(err)
	}
	return fb
}

func createSoundBuffer(srcPath string) utils.FrameBuffer {
	fb := &utils.IpcPipelineBuffer{}
	err := fb.Open(srcPath, 8*1024, 1024*1024)
	if err != nil {
		panic(err)
	}
	return fb
}

func run(cmd *cobra.Command, args []string) {
	wh := strings.Split(resolution, "x")
	width, _ := strconv.Atoi(wh[0])
	height, _ := strconv.Atoi(wh[1])
	fps, _ := strconv.Atoi(fps)

	fi := renderer.CreateFrameInfo(width, height, fps)
	r := renderer.CreateRenderer(renderer.TypeWebRTCStreamRenderer, fi)

	ih := input.CreateInputHandler(keyEvtQueue)
	fb := createFrameBuffer(frameSource, fi)
	sb := createSoundBuffer(soundSource)
	defer fb.Close()

	r.StartWithFrameBuffer(fb, sb, ih)
}

func _main() {
	entryCmd := cmdRunApp()
	if err := entryCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
