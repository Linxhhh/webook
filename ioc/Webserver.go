package ioc

import (
	"github.com/Linxhhh/webook/internal/app"
	"github.com/gin-gonic/gin"
)

func InitWebServer(halFunc []gin.HandlerFunc, userHdl *app.UserHandler) *gin.Engine {
	router := gin.Default()
	router.Use(halFunc...)
	userHdl.RegistryRouter(router)
	return router
}