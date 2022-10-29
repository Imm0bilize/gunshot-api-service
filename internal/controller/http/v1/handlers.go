package v1

import (
	"github.com/Imm0bilize/gunshot-api-service/internal/controller/http/dto"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository"
	"github.com/Imm0bilize/gunshot-api-service/internal/uCase"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) UploadAudio(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (h *Handler) RegisterNewClient(c *gin.Context) {
	var clientRequest dto.NewClientRequest

	if err := c.ShouldBindJSON(&clientRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info(
		"got new request for creating new clientRequest",
		zap.String("reqID", clientRequest.RequestID.String()),
	)

	id, err := h.domain.CreateNewClient(
		c.Request.Context(), clientRequest.RequestID.String(), clientRequest.ClientInformation,
	)

	if err != nil {
		if errors.Is(err, uCase.ErrRequestAlreadyProcessed) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"clientID": id})
	return
}

func (h *Handler) UpdateClient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (h *Handler) DeleteClient(c *gin.Context) {
	var request dto.DeleteClientRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err := h.domain.DeleteClient(c.Request.Context(), request.RequestID.String(), request.ClientID.String())
	if err != nil {
		if errors.As(err, &repository.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "the client not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "ok")
}
