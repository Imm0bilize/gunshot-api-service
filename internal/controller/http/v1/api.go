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

func (h *Handler) InitAPI(router *gin.RouterGroup,
	injectRequestID, injectClientID func(c *gin.Context),
) {

	v1 := router.Group("v1")
	{
		client := v1.Group("client")
		{
			client.Use(injectRequestID)

			client.POST("", h.RegisterNewClient)

			clientID := client.Group(":id")
			{
				clientID.Use(injectClientID)

				clientID.GET("", h.GetClient)
				clientID.PUT("", h.UpdateClient)
				clientID.DELETE("", h.DeleteClient)

				clientID.POST(":ts/upload", h.UploadAudio)
			}
		}
	}
}
