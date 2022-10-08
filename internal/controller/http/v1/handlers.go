package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) UploadAudio(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (h *Handler) RegisterNewClient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (h *Handler) DeleteClient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
