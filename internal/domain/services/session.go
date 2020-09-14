package services

import (
	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/models"
)

type WebRTCSession interface {
	ProcessNewOffer(
		sdp models.SdpInfo) (string, error)

	StartIceProcess(
		peerId string,
		onLocalIceCandidate func(b64EncodedIceCandidate string),
		onIceStateChanged func(connectionState string)) error

	ProcessRemoteIce(
		iceCandidate models.IceCandidate) error

	GetRenderContext() engine.RenderContext
	GetInputContext() engine.InputContext
}
