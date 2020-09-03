package services

import "github.com/oraksil/orakki/internal/domain/engine"

type WebRTCSession interface {
	OnCreatedOffer(
		peerId string,
		b64EncodedOffer string,
		onCreatedAnswer func(peerId, b64EncodedAnswer string)) error

	StartIceProcess(
		peerId string,
		onLocalIceCandidates func(peerId, b64EncodedIceCandidate string),
		onIceStateChanged func(peerId, connectionState string)) error

	OnRemoteIceCandidates(
		peerId string,
		b64EncodedIceCandidate string) error

	GetRenderContext() engine.RenderContext
	GetInputContext() engine.InputContext
}
