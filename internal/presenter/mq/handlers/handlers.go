package handlers

import (
	"fmt"
	"time"

	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/usecases"
	"github.com/oraksil/orakki/internal/presenter/mq/dto"
	"github.com/sangwonl/mqrpc"
	// "time"
)

type SystemHandler struct {
	SystemMonitorUseCase *usecases.SystemStateMonitorUseCase
}

func (h *SystemHandler) handleFetchState(ctx *mqrpc.Context) interface{} {
	sysState, _ := h.SystemMonitorUseCase.GetSystemState()
	return dto.SystemStateToOrakkiState(sysState)
}

func (h *SystemHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_FETCH_ORAKKI_STATE, Handler: h.handleFetchState},
	}
}

type SetupHandler struct {
	SetupUseCase *usecases.SetupUseCase
}

func (h *SetupHandler) handleSetupOffer(ctx *mqrpc.Context) interface{} {
	return &dto.SetupAnswer{PeerId: "some id", Answer: "Some Answer"}
}

type Icecandidate struct {
	PlayerId  int64
	OrakkiId  string
	IceString string
}

func (h *SetupHandler) addIceCandidate(ctx *mqrpc.Context) interface{} {
	h.SetupUseCase.MessageService.Request("test_peer", models.MSG_HANDLE_SETUP_ICECANDIDATE, Icecandidate{OrakkiId: "orakki1", IceString: "Server Ice Candidate1"}, 5*time.Second)
	h.SetupUseCase.MessageService.Request("test_peer", models.MSG_HANDLE_SETUP_ICECANDIDATE, Icecandidate{OrakkiId: "orakki1", IceString: "Server Ice Candidate2"}, 5*time.Second)
	h.SetupUseCase.MessageService.Request("test_peer", models.MSG_HANDLE_SETUP_ICECANDIDATE, Icecandidate{OrakkiId: "orakki1", IceString: "Server Ice Candidate3"}, 5*time.Second)
	h.SetupUseCase.MessageService.Request("test_peer", models.MSG_HANDLE_SETUP_ICECANDIDATE, Icecandidate{OrakkiId: "orakki1", IceString: ""}, 5*time.Second)
	return &dto.SetupAnswer{PeerId: "some id", Answer: "Some Answer"}
}

func (h *SetupHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_HANDLE_SETUP_OFFER, Handler: h.handleSetupOffer},
		{MsgType: models.MSG_HANDLE_SETUP_ICECANDIDATE, Handler: h.addIceCandidate},
	}
}
