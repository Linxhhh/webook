// +build wireinject

package main

import (
	"github.com/Linxhhh/webook/internal/app"
	"github.com/Linxhhh/webook/internal/repository"
	"github.com/Linxhhh/webook/internal/repository/cache"
	"github.com/Linxhhh/webook/internal/repository/dao"
	"github.com/Linxhhh/webook/internal/service"
	"github.com/Linxhhh/webook/internal/events"
	"github.com/Linxhhh/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitCache, ioc.InitDB, ioc.InitSmsService, ioc.InitSaramaClient, ioc.InitSyncProducer,

		// DAO
		dao.NewUserDAO,
		dao.NewArticleDAO,
		dao.NewInteractionDAO,
		dao.NewFollowDAO,
		dao.NewFeedPushEventDAO,
		dao.NewFeedPullEventDAO,

		// Cache
		cache.NewUserCache,
		cache.NewCodeCache,
		cache.NewArticleCache,
		cache.NewInteractionCache,
		cache.NewFollowCache,
		cache.NewFeedEventCache,

		// Repository
		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewArticleRepository,
		repository.NewInteractionRepository,
		repository.NewFollowRepository,
		repository.NewFeedEventRepo,

		// Service
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		service.NewInteractionService,
		service.NewFollowService,
		service.NewFeedEventService,

		// Event
		events.NewArticleEventProducer,
		events.NewArticleEventConsumer,
		ioc.InitConsumers,

		// Handler
		app.NewUserHandler,
		app.NewArticleHandler,
		app.NewFollowHandler,

		// Webserver
		ioc.InitMiddleware,
		ioc.InitWebServer,
	)

	return gin.Default()
}