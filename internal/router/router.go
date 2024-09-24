package router

import (
	"github.com/MHSaeedkia/tinyTimer/internal/handler"
	"github.com/gin-gonic/gin"
)

func GetEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	e := gin.Default()
	e.GET("/start", handler.Start)
	e.GET("/stop", handler.Stop)
	e.GET("/chart", handler.TotalTime)
	e.GET("/clear", handler.Clear)
	return e
}
