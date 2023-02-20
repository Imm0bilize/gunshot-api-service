package v1

import (
	"fmt"
	"github.com/Imm0bilize/gunshot-api-service/internal/controller/http/dto"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository"
	"github.com/Imm0bilize/gunshot-api-service/internal/uCase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) RegisterNewClient(c *gin.Context) {
	var (
		req       dto.ClientInfo
		requestID = c.MustGet("requestID").(uuid.UUID)
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Msg: err.Error()})
		return
	}

	h.logger.Info(
		"got new request for creating new client",
		zap.String("request_id", requestID.String()),
	)

	id, err := h.domain.Client.Create(
		c.Request.Context(),
		requestID,
		&entities.Client{
			FullName:            req.FullName,
			LocationName:        req.LocationName,
			Latitude:            req.Latitude,
			Longitude:           req.Longitude,
			NotificationMethods: req.NotificationMethods,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{ClientID: id})
	return
}

func (h *Handler) UpdateClient(c *gin.Context) {
	var (
		requestID = c.MustGet("requestID").(uuid.UUID)
		clientID  = c.MustGet("clientID").(string)
		req       dto.ClientInfo
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Msg: err.Error()})
		return
	}

	err := h.domain.Client.Update(
		c.Request.Context(),
		requestID,
		clientID,
		&entities.Client{
			LocationName:        req.LocationName,
			FullName:            req.FullName,
			Latitude:            req.Latitude,
			Longitude:           req.Longitude,
			NotificationMethods: req.NotificationMethods,
		},
	)

	if err != nil {
		if errors.Is(err, repository.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Msg: err.Error()})
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Msg: err.Error()})
	}

}

func (h *Handler) GetClient(c *gin.Context) {
	var (
		requestID = c.MustGet("requestID").(uuid.UUID)
		clientID  = c.MustGet("clientID").(string)
	)

	client, err := h.domain.Client.Get(c.Request.Context(), requestID, clientID)
	if err != nil {
		if errors.Is(err, repository.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Msg: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *Handler) DeleteClient(c *gin.Context) {
	var (
		clientID  = c.MustGet("clientID").(string)
		requestID = c.MustGet("requestID").(uuid.UUID)
	)

	if err := h.domain.Client.Delete(c.Request.Context(), requestID, clientID); err != nil {
		if errors.Is(err, repository.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "the client not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "ok")
}

func (h *Handler) UploadAudio(c *gin.Context) {
	var (
		requestID = c.MustGet("requestID").(uuid.UUID)
		clientID  = c.MustGet("clientID").(string)
		timestamp time.Time
		payload   []byte
	)

	nsec, err := strconv.Atoi(c.Param("ts"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("can't convert string nsec to int: %s", err.Error()),
		})

		return
	}

	timestamp = time.Unix(0, int64(nsec))

	fileHeader, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("can't get form file (audio): %s", err.Error()),
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("can't open sended file: %s", err.Error()),
		})
		return
	}

	payload, err = ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("can't read sended file: %s", err.Error()),
		})
		return
	}

	err = h.domain.Audio.Upload(c.Request.Context(), requestID, clientID,
		entities.AudioMessage{
			Timestamp: timestamp,
			Payload:   payload,
		})

	if err != nil {
		if errors.Is(err, uCase.ErrNotEqRequiredLength) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Msg: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("can't send audio for processing: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}
