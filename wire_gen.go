package main

import (
	"github.com/Linxhhh/webook/internal/app"
	"github.com/Linxhhh/webook/internal/events"
	"github.com/Linxhhh/webook/internal/repository"
	"github.com/Linxhhh/webook/internal/repository/cache"
	"github.com/Linxhhh/webook/internal/repository/dao"
	"github.com/Linxhhh/webook/internal/service"
	"github.com/Linxhhh/webook/ioc"
	"github.com/gin-gonic/gin"
)

type WebServer struct {
	engine    *gin.Engine
	consumers []events.Consumer
}

func InitWebServer() WebServer {

	// 第三方依赖
	m, s := ioc.InitDB()
	cmdable := ioc.InitCache()
	smsService := ioc.InitSmsService()
	sclient := ioc.InitSaramaClient()
	sproducer := ioc.InitSyncProducer(sclient)

	// DAO
	userDAO := dao.NewUserDAO(m, s)
	articleDAO := dao.NewArticleDAO(m, s)
	interactionDAO := dao.NewInteractionDAO(m, s)
	followDAO := dao.NewFollowDAO(m, s)
	feedPullEventDAO := dao.NewFeedPullEventDAO(m, s)
	feedPushEventDAO := dao.NewFeedPushEventDAO(m, s)

	// Cache
	userCache := cache.NewUserCache(cmdable)
	codeCache := cache.NewCodeCache(cmdable)
	articleCache := cache.NewArticleCache(cmdable)
	interactionCache := cache.NewInteractionCache(cmdable)
	followCache := cache.NewFollowCache(cmdable)
	feedEventCache := cache.NewFeedEventCache(cmdable)

	// Repository
	userRepository := repository.NewUserRepository(userDAO, userCache)
	codeRepository := repository.NewCodeRepository(codeCache)
	articleRepository := repository.NewArticleRepository(articleDAO, articleCache)
	interactionRepository := repository.NewInteractionRepository(interactionDAO, interactionCache)
	followRepository := repository.NewFollowRepository(followDAO, followCache)
	feedEventRepository := repository.NewFeedEventRepo(feedPullEventDAO, feedPushEventDAO, feedEventCache)

	// Service
	userService := service.NewUserService(userRepository)
	codeService := service.NewCodeService(codeRepository, smsService)
	articleService := service.NewArticleService(articleRepository, userRepository)
	interactionService := service.NewInteractionService(interactionRepository, articleRepository)
	followService := service.NewFollowService(followRepository)
	feedEventService := service.NewFeedEventService(feedEventRepository, followRepository)

	// Event
	articleEventProducer := events.NewArticleEventProducer(sproducer)
	articleEventConsumer := events.NewArticleEventConsumer(sclient, feedEventService)

	// Handler
	userHandler := app.NewUserHandler(userService, codeService)
	articleHandler := app.NewArticleHandler(articleService, interactionService, articleEventProducer)
	followHandler := app.NewFollowHandler(followService)

	// Webserver
	v := ioc.InitMiddleware()
	engine := ioc.InitEngine(v, userHandler, articleHandler, followHandler)
	consumers := ioc.InitConsumers(articleEventConsumer)
	
	return WebServer{
		engine: engine,
		consumers: consumers,
	}
}
