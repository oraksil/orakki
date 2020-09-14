package usecases

import (
	"github.com/oraksil/orakki/internal/domain/models"
	"github.com/oraksil/orakki/internal/domain/services"
)

type SystemUseCase struct {
	ServiceConfig  *services.ServiceConfig
	MessageService services.MessageService
}

func (uc *SystemUseCase) GetOrakkiState() (*models.OrakkiState, error) {
	return &models.OrakkiState{
		OrakkiId: uc.ServiceConfig.OrakkiId,
		State:    models.ORAKKI_STATE_READY,
	}, nil
}
