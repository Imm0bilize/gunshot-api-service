package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const _requestIDHeader = "X-REQUEST-ID"

func InjectRequestIDIntoCtx(c *gin.Context) {
	var requestID uuid.UUID

	if err := requestID.Scan(c.GetHeader(_requestIDHeader)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid or empty header 'X-REQUEST-ID'"})
		return
	}

	c.Set("requestID", requestID)
	c.Next()
}

func InjectClientIDIntoCtx(c *gin.Context) {
	var clientID string

	if clientID = c.Param("id"); clientID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid clientID: clientID must not be empty"})
		return
	}

	c.Set("clientID", clientID)
	c.Next()
}
