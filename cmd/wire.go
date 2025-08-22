//go:build wireinject
// +build wireinject

package main

import (
	"Nuxus/configs"
	"Nuxus/internal/controller"
	"Nuxus/internal/dao"
	"Nuxus/internal/middleware"
	"Nuxus/internal/routers"
	"Nuxus/internal/service"
	"Nuxus/internal/tasks"
	"github.com/google/wire"
)

type App struct {
	Router              *routers.Router
	SyncTask            *tasks.SyncTask
	Config              *configs.Config
	MiddlewareManager   *middleware.MiddlewareManager
}

func NewApp(
	router *routers.Router,
	syncTask *tasks.SyncTask,
	config *configs.Config,
	middlewareManager *middleware.MiddlewareManager,
) *App {
	return &App{
		Router:            router,
		SyncTask:          syncTask,
		Config:            config,
		MiddlewareManager: middlewareManager,
	}
}

// Wire Provider Set
var ProviderSet = wire.NewSet(
	// Config
	configs.LoadConfig,
	
	// 基础设施层
	dao.NewDB,
	dao.NewClient,
	dao.NewRedisClient,
	dao.NewRepository,
	
	// DAO层
	dao.NewUserDAO,
	dao.NewPostDAO,
	dao.NewTagDAO,
	
	// Middleware层
	middleware.NewMiddlewareManager,
	
	// Service层
	service.NewEmailService,
	service.NewAccountService,
	service.NewUserService,
	service.NewPostService,
	service.NewTagService,
	
	// Controller层
	controller.NewUserController,
	controller.NewPostController,
	controller.NewTagController,
	
	// Router层
	routers.NewRouter,
	
	// Tasks
	tasks.NewSyncTask,
	
	// App
	NewApp,
)

func InitializeApp() (*App, error) {
	wire.Build(ProviderSet)
	return &App{}, nil
}