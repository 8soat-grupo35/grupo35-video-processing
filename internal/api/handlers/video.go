package handlers

import (
	"github.com/8soat-grupo35/fastfood-order-production/internal/interfaces/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type VideoHandler struct {
	videoUploaderUseCase usecase.VideoUploaderUseCase
}

func NewVideoHandler(useCase usecase.VideoUploaderUseCase) VideoHandler {
	return VideoHandler{
		videoUploaderUseCase: useCase,
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
	files, err := c.MultipartForm()

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	videoFiles := files.File["videos"]

	// Call the use case
	err = v.videoUploaderUseCase.UploadVideos(videoFiles)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Videos uploaded successfully")
}
