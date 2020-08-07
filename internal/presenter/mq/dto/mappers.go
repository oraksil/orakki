package dto

import "gitlab.com/oraksil/orakki/internal/domain/models"

func SystemStateToOrakkiState(src *models.SystemState) *OrakkiState {
	return &OrakkiState{
		OrakkiId: src.OrakkiId,
		State:    src.OrakkiState,
	}
}
