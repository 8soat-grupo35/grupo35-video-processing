package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/usecases"
	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	transferVideoUseCase usecases.TransferVideoUseCase
	processVideoSender   usecases.ProcessVideoSender
	videoRepository      usecases.VideoRepository
}

func NewVideoHandler(transferUseCase usecases.TransferVideoUseCase, processVideoSender usecases.ProcessVideoSender, videoRepository usecases.VideoRepository) VideoHandler {
	return VideoHandler{
		transferVideoUseCase: transferUseCase,
		processVideoSender:   processVideoSender,
		videoRepository:      videoRepository,
	}
}

// ListByUser handles the request to list videos by a specific user.
//
// @Summary List videos by user
// @Description Retrieves a list of videos associated with a specific user ID.
// @Tags Videos
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} Video "List of videos"
// @Failure 400 {string} string "userID is required"
// @Failure 404 {string} string "No videos found for this user"
// @Failure 500 {string} string "Internal server error"
// @Router /videos/user/{userId} [get]
func (v *VideoHandler) ListByUser(c echo.Context) error {
	userID := c.Param("userId")

	if userID == "" {
		return c.JSON(http.StatusBadRequest, "userID is required")
	}

	videos, err := v.videoRepository.GetByUserId(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(videos) == 0 {
		return c.JSON(http.StatusNotFound, "No videos found for this user")
	}

	return c.JSON(http.StatusOK, videos)
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
	videoFiles, user, err := v.validateAndExtractFiles(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Processa cada arquivo de vídeo
	for _, file := range videoFiles {
		if err := v.processVideoFile(c, file, *user); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, "Videos uploaded successfully")
}

// validateAndExtractFiles valida os arquivos do formulário e extrai os vídeos e o userID.
func (v *VideoHandler) validateAndExtractFiles(c echo.Context) ([]*multipart.FileHeader, *entities.User, error) {
	files, err := c.MultipartForm()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse form: %w", err)
	}

	userID := c.Get("user_id")
	userEmail := c.Get("user_email")

	if userID == nil || userEmail == nil {
		return nil, nil, fmt.Errorf("user ID or email not found in context")
	}

	user := &entities.User{
		ID:    userID.(string),
		Email: userEmail.(string),
	}

	videoFiles := files.File["videos"]
	if len(videoFiles) == 0 {
		return nil, nil, fmt.Errorf("no video files provided")
	}

	return videoFiles, user, nil
}

// processVideoFile processa um único arquivo de vídeo.
func (v *VideoHandler) processVideoFile(c echo.Context, file *multipart.FileHeader, user entities.User) error {
	// Abre o arquivo
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Lê os dados do arquivo
	videoData, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Gera a chave para o arquivo no S3
	key := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	contentType := http.DetectContentType(videoData)
	// Faz o upload do vídeo para o S3
	if err := v.transferVideoUseCase.UploadVideo(key, videoData, contentType); err != nil {
		fmt.Errorf("failed to upload video to S3: %w", err)
	}

	// Cria mensagem para o SQS
	videoMessage := entities.VideoMessage{
		User:      user,
		VideoPath: key,
	}

	// Envia o vídeo para a fila de processamento
	if err := v.processVideoSender.Send(videoMessage); err != nil {
		return fmt.Errorf("failed to send video to processing queue: %w", err)
	}

	// Cria a entidade do vídeo
	video := entities.Video{
		UserID:    user.ID,
		Path:      key,
		Status:    entities.VideoStatusInProcessing,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Salva os metadados do vídeo no DynamoDB
	if _, err := v.videoRepository.Create(video); err != nil {
		return fmt.Errorf("failed to save video metadata: %w", err)
	}

	// Loga o upload bem-sucedido
	fmt.Printf("Uploaded video: %s\n", key)
	return nil
}
