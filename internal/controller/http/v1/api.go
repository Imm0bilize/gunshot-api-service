package v1

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	tracer trace.Tracer
	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) InitApi(router *gin.RouterGroup) {
	v1 := router.Group("v1")
	{
		v1.POST("/register", h.RegisterNewClient)
		user := v1.Group("/:id")
		{
			user.DELETE("", h.DeleteClient)
			user.POST(":ts/upload", h.UploadAudio)
		}
	}
}
