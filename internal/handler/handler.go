package handler

import (
	"fmt"
	"time"

	"github.com/MHSaeedkia/tinyTimer/internal/db"
	"github.com/gin-gonic/gin"
)

func Start(ctx *gin.Context) {
	dD := db.GetDB()
	db.Start(dD, time.Now())
}

func Stop(ctx *gin.Context) {
	dD := db.GetDB()
	db.Stop(dD)
}

func TotalTime(ctx *gin.Context) {
	dD := db.GetDB()
	_, result := db.TotalTime(dD)
	fmt.Println(result)
}

func Clear(ctx *gin.Context) {
	dD := db.GetDB()
	db.Clear(dD)
}
