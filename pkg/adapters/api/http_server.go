package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KScaesar/GARMIN2024/pkg"
	"github.com/KScaesar/GARMIN2024/pkg/app"
)

func NewHttpServer(conf *pkg.Config, mux http.Handler) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.HttpPort),
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func NewGinRouter(conf *pkg.Config, svc *app.Service) *gin.Engine {
	router := gin.New()
	if !conf.GinDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	router.Use(gin.Recovery())

	router.Group("api/v1/").
		POST("orders", createOrder(svc))

	return router
}
