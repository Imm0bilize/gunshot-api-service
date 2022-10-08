package http

import (
	v1 "github.com/Imm0bilize/gunshot-api-service/internal/controller/http/v1"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
)

func New(logger *zap.Logger, middlewares ...gin.HandlerFunc) (*gin.Engine, error) {
	router := gin.New()

	router.Use(middlewares...)

	//trace
	router.Use(otelgin.Middleware("gunshot-api-service"))

	// Debug handlers
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	initPprof(router.Group("/debug"))

	// API
	initApi(router, logger)

	return router, nil
}

func initPprof(router *gin.RouterGroup) {
	p := router.Group("/pprof")
	{
		p.GET("/", gin.WrapF(pprof.Index))
		p.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		p.GET("/profile", gin.WrapF(pprof.Profile))
		p.POST("/symbol", gin.WrapF(pprof.Symbol))
		p.GET("/symbol", gin.WrapF(pprof.Symbol))
		p.GET("/trace", gin.WrapF(pprof.Trace))
		p.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		p.GET("/block", gin.WrapH(pprof.Handler("block")))
		p.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		p.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		p.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		p.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}

func initApi(router *gin.Engine, logger *zap.Logger) {
	handlerV1 := v1.NewHandler(logger)

	api := router.Group("/api")
	{
		handlerV1.InitApi(api)
	}
}
