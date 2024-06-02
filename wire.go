// +build wireinject

package main

import (
	"github.com/Linxhhh/webook/internal/app"
	"github.com/Linxhhh/webook/internal/repository"
	"github.com/Linxhhh/webook/internal/repository/cache"
	"github.com/Linxhhh/webook/internal/repository/dao"
	"github.com/Linxhhh/webook/internal/service"
	"github.com/Linxhhh/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitCache, ioc.InitDB, ioc.InitSmsService,

		// DAO
		dao.NewUserDAO,

		// Cache
		cache.NewUserCache,
		cache.NewCodeCache,

		// Repository
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// Service
		service.NewUserService,
		service.NewCodeService,

		// Handler
		app.NewUserHandler,

		// Webserver
		ioc.InitMiddleware,
		ioc.InitWebServer,
	)

	return gin.Default()
}