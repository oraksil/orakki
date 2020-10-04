package services

import (
	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/models"
)

type WebRTCSession interface {
	SetupIceHandlers(
		peerInfo models.PeerInfo,
		onLocalIceCandidate func(iceCandidate models.IceCandidate),
		onIceStateChanged func(peerInfo models.PeerInfo, connectionState string)) error
	ProcessNewOffer(sdp models.SdpInfo) (*models.SdpInfo, error)
	ProcessRemoteIce(iceCandidate models.IceCandidate) error

	GetRenderContext() engine.RenderContext
	GetInputContext() engine.InputContext
	GetSessionContext() engine.SessionContext
}
