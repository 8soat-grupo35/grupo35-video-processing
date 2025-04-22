package handlers

import (
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/usecases"
)

type StatusHandler struct {
	statusVideoConsumer usecases.StatusVideoConsumer
}

func NewStatusHandler(statusVideoConsumer usecases.StatusVideoConsumer) StatusHandler {
	return StatusHandler{
		statusVideoConsumer: statusVideoConsumer,
	}
}

func (s StatusHandler) UpdateStatus() {
	s.statusVideoConsumer.ConsumeMessages()
}
