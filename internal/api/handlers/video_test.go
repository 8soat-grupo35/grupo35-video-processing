package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	mock_usecases "github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/usecases/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func createMultipartForm(t *testing.T, fields map[string]string, files map[string][]byte) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}
	for fname, content := range files {
		part, err := writer.CreateFormFile("videos", fname)
		assert.NoError(t, err)
		_, err = part.Write(content)
		assert.NoError(t, err)
	}
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestUpload_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransfer := mock_usecases.NewMockTransferVideoUseCase(ctrl)
	mockSender := mock_usecases.NewMockProcessVideoSender(ctrl)
	mockRepo := mock_usecases.NewMockVideoRepository(ctrl)

	handler := NewVideoHandler(mockTransfer, mockSender, mockRepo)

	videoContent := []byte("video-bytes")
	fields := map[string]string{
		"userID":    "user-1",
		"userEmail": "user@email.com",
	}
	files := map[string][]byte{
		"test.mp4": videoContent,
	}
	body, contentType := createMultipartForm(t, fields, files)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/videos/upload", body)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Espera chamadas dos mocks
	mockTransfer.EXPECT().UploadVideo(gomock.Any(), videoContent, gomock.Any()).Return(nil)
	mockSender.EXPECT().Send(gomock.Any()).Return(nil)
	mockRepo.EXPECT().Create(gomock.Any()).Return(&entities.Video{}, nil)

	err := handler.Upload(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Videos uploaded successfully")
}

func TestUpload_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := NewVideoHandler(nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/videos/upload", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Upload(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProcessVideoFile_ErrorOnOpen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransfer := mock_usecases.NewMockTransferVideoUseCase(ctrl)
	mockSender := mock_usecases.NewMockProcessVideoSender(ctrl)
	mockRepo := mock_usecases.NewMockVideoRepository(ctrl)

	handler := NewVideoHandler(mockTransfer, mockSender, mockRepo)

	file := &multipart.FileHeader{
		Filename: "test.mp4",
		Size:     0,
	}
	user := entities.User{ID: "user-1", Email: "user@email.com"}

	// Simula erro ao abrir o arquivo
	err := handler.processVideoFile(nil, file, user)
	assert.Error(t, err)
}
