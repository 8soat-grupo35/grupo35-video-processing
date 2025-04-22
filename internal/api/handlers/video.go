package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/usecases"
	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	transferVideoUseCase usecases.TransferVideoUseCase
	videoProcessSender   usecases.VideoProcessSender
	videoRepository      usecases.VideoRepository
}

func NewVideoHandler(transferUseCase usecases.TransferVideoUseCase, videoProcessSender usecases.VideoProcessSender, videoRepository usecases.VideoRepository) VideoHandler {
	return VideoHandler{
		transferVideoUseCase: transferUseCase,
		videoProcessSender:   videoProcessSender,
		videoRepository:      videoRepository,
	}
}

// Upload godoc
// @Summary Upload videos
// @Description Upload multiple video files
// @Tags videos
// @Accept multipart/form-data
// @Produce json
// @Param videos formData file true "Video files"
// @Success 200 {string} string "Videos uploaded successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/videos/upload [post]
func (v *VideoHandler) Upload(c echo.Context) error {
	// Valida e obtém os arquivos do formulário
	videoFiles, userID, err := v.validateAndExtractFiles(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Processa cada arquivo de vídeo
	for _, file := range videoFiles {
		if err := v.processVideoFile(c, file, userID); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, "Videos uploaded successfully")
}

// validateAndExtractFiles valida os arquivos do formulário e extrai os vídeos e o userID.
func (v *VideoHandler) validateAndExtractFiles(c echo.Context) ([]*multipart.FileHeader, string, error) {
	files, err := c.MultipartForm()
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse form: %w", err)
	}

	userID := c.FormValue("userID")
	if userID == "" {
		return nil, "", fmt.Errorf("userID is required")
	}

	videoFiles := files.File["videos"]
	if len(videoFiles) == 0 {
		return nil, "", fmt.Errorf("no video files provided")
	}

	return videoFiles, userID, nil
}

// processVideoFile processa um único arquivo de vídeo.
func (v *VideoHandler) processVideoFile(c echo.Context, file *multipart.FileHeader, userID string) error {
	// Abre o arquivo
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Lê os dados do arquivo
	videoData := make([]byte, file.Size)
	if _, err := src.Read(videoData); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Gera a chave para o arquivo no S3
	key := fmt.Sprintf("videos/%d_%s", time.Now().UnixNano(), file.Filename)

	// Faz o upload do vídeo para o S3
	if err := v.transferVideoUseCase.UploadVideo(key, videoData); err != nil {
		return fmt.Errorf("failed to upload video to S3: %w", err)
	}

	// Cria a entidade do vídeo
	video := entities.Video{
		UserID:    userID,
		Path:      key,
		Status:    entities.VideoStatusInProcessing,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Envia o vídeo para a fila de processamento
	if err := v.videoProcessSender.Send(video); err != nil {
		return fmt.Errorf("failed to send video to processing queue: %w", err)
	}

	// Salva os metadados do vídeo no DynamoDB
	if _, err := v.videoRepository.Create(video); err != nil {
		return fmt.Errorf("failed to save video metadata: %w", err)
	}

	// Loga o upload bem-sucedido
	fmt.Printf("Uploaded video: %s\n", key)
	return nil
}
