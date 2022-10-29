package v1

import (
	"github.com/Imm0bilize/gunshot-api-service/internal/uCase"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	tracer trace.Tracer
	logger *zap.Logger
	domain *uCase.UseCase
}

func NewHandler(logger *zap.Logger, domain *uCase.UseCase) *Handler {
	return &Handler{
		logger: logger,
		domain: domain,
	}
}

func (h *Handler) InitAPI(router *gin.RouterGroup) {
	v1 := router.Group("v1")
	{
		v1.POST("/register", h.RegisterNewClient)

		client := v1.Group("client")
		{
			client.DELETE("", h.DeleteClient)
			client.PUT("", h.UpdateClient)
			client.POST("/upload", h.UploadAudio)
		}
	}
}
