package handlers

import (
	"testing"

	mock_usecases "github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/usecases/mock"
	"go.uber.org/mock/gomock"
)

func TestStatusHandler_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConsumer := mock_usecases.NewMockStatusVideoConsumer(ctrl)
	mockConsumer.EXPECT().ConsumeMessages().Times(1)

	handler := NewStatusHandler(mockConsumer)
	handler.UpdateStatus()
}
