package usecases

import (
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
)

type SystemStateMonitorUseCase struct {
	ServiceConfig  *services.ServiceConfig
	MessageService services.MessageService
}

func (uc *SystemStateMonitorUseCase) GetSystemState() (*models.SystemState, error) {
	return &models.SystemState{
		OrakkiId:    uc.ServiceConfig.OrakkiId,
		OrakkiState: models.ORAKKI_STATE_READY,
	}, nil
}
