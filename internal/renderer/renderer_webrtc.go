package renderer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/oraksil/orakki/internal/input"
	"github.com/oraksil/orakki/pkg/utils"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

type WrRenderer struct {
	router  *gin.Engine
	fbImage utils.FrameBuffer
	fbSound utils.FrameBuffer
	fi      FrameInfo
	ih      *input.InputHandler

	upgrader websocket.Upgrader

	videoTrack *webrtc.Track
	audioTrack *webrtc.Track
}

type Msg struct {
	Datatype string
	Message  string
}

func JsonDeserialize(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func JsonSerialize(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

func (r *WrRenderer) StartWithFrameBuffer(fb utils.FrameBuffer, sb utils.FrameBuffer, ih *input.InputHandler) {
	r.fbImage = fb
	r.fbSound = sb
	r.ih = ih

	inputReader := input.CreateReader(100)
	r.ih.SetInputReader(inputReader)
	r.ih.Run()

	go r.videoWriter()
	go r.audioWriter()

	r.router.Run(":8000")
}

func (r *WrRenderer) routerHandlePage(c *gin.Context) {
	c.HTML(http.StatusOK, "renderer_webrtc.html", gin.H{})
}

func (r *WrRenderer) routerHandleMeta(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"meta": &FrameMetaDTO{
		Fps: r.fi.fps,
		W:   r.fi.w,
		H:   r.fi.h,
	}})
}

func (r *WrRenderer) videoWriter() {
	for {
		select {
		default:
			buf := r.fbImage.GetBuffer()
			if buf != nil {
				r.videoTrack.WriteSample(media.Sample{Data: buf, Samples: 1})
			}
		}
	}
}

func (r *WrRenderer) audioWriter() {
	for {
		select {
		default:
			buf := r.fbSound.GetBuffer()
			if buf != nil {
				go func() {
					r.audioTrack.WriteSample(media.Sample{Data: buf, Samples: 960})
				}()
			}
		}
	}
}

func (r *WrRenderer) routerHandleSignaling(c *gin.Context) {
	ws, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			fmt.Printf("=== Ice Candidate : %s\n", c)
			a := c.ToJSON()

			err = ws.WriteJSON(Msg{
				Datatype: "icecandidate",
				Message:  JsonSerialize(a)})

			if err != nil {
				ws.Close()
				peerConnection.Close()
			}
		}
	})
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("=== ICE Connection State has changed: %s\n", connectionState.String())
	})

	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(r.videoTrack)
	_, err = peerConnection.AddTrack(r.audioTrack)
	if err != nil {
		panic(err)
	}

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("=== New DataChannel %s %d\n", d.Label(), d.ID())

		d.OnOpen(func() {
			r.ih.SetPadInOrder(d.ID())
		})

		d.OnClose(func() {
			r.ih.RemovePad(d.ID())
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			gipanKey := input.ConvertToGipanKeys(r.ih.GetPlayerOrderFromPad(d.ID()), msg.Data)
			r.ih.Reader.AddInput(gipanKey)
		})
	})

	if err != nil {
		panic(err)
	}

	for {
		m := Msg{}
		err := ws.ReadJSON(&m)

		if err != nil {
			ws.Close()
			peerConnection.Close()
			break
		}

		if m.Datatype == "icecandidate" {
			candidate := webrtc.ICECandidateInit{}
			JsonDeserialize(m.Message, &candidate)
			peerConnection.AddICECandidate(candidate)
			fmt.Printf("=== Ice Candidate arrived : %s\n", candidate.Candidate)
		} else if m.Datatype == "sdp" {
			offer := webrtc.SessionDescription{}
			JsonDeserialize(m.Message, &offer)
			fmt.Printf("=== Offer arriver : %s\n", offer)

			err = peerConnection.SetRemoteDescription(offer)
			if err != nil {
				ws.Close()
				peerConnection.Close()
				break
			}

			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				ws.Close()
				peerConnection.Close()
				break
			}

			err = ws.WriteJSON(Msg{
				Datatype: "sdp",
				Message:  JsonSerialize(answer)})

			if err != nil {
				ws.Close()
				peerConnection.Close()
				break
			}

			err = peerConnection.SetLocalDescription(answer)
			if err != nil {
				ws.Close()
				peerConnection.Close()
				break
			}
		}
	}
}

func createWebRTCStreamRenderer(frameInfo FrameInfo) Renderer {
	r := WrRenderer{
		router:   gin.Default(),
		fbImage:  nil,
		fbSound:  nil,
		fi:       frameInfo,
		upgrader: websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
	}

	codec := webrtc.NewRTPH264Codec(webrtc.DefaultPayloadTypeH264, 8000)
	videoTrack, err := webrtc.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "video", "pion2", codec)
	if err != nil {
		panic(err)
	}

	codec = webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 8000)
	audioTrack, err := webrtc.NewTrack(webrtc.DefaultPayloadTypeOpus, rand.Uint32(), "audio", "pion2", codec)
	if err != nil {
		panic(err)
	}

	r.videoTrack = videoTrack
	r.audioTrack = audioTrack
	r.router.Use(cors.Default())
	r.router.LoadHTMLFiles("web/templates/renderer_webrtc.html")
	r.router.GET("/", r.routerHandlePage)
	r.router.GET("/meta", r.routerHandleMeta)
	r.router.GET("/signaling", r.routerHandleSignaling)
	r.router.Static("/js", "web/js")
	r.router.Static("/css", "web/css")

	return &r
}
