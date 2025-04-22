package usecases

import (
	"testing"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	repository_mocks "github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository/mock"
	"go.uber.org/mock/gomock"
)

func TestProcessVideoSender_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSQS := repository_mocks.NewMockSQS(ctrl)
	videoMessage := entities.VideoMessage{
		User: entities.User{
			ID:    "user-id",
			Email: "user-email",
		},
		VideoPath: "video-path",
	}

	mockSQS.EXPECT().SendMessage(videoMessage).Return(nil)

	sender := NewProcessVideoSender(mockSQS)
	err := sender.Send(videoMessage)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
