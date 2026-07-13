package app

import (
	"log"
	"project-shopping-cart/internal/config"
	"project-shopping-cart/internal/routes"
	"project-shopping-cart/internal/validation"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
	modules []Module
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validator init failed %v", err)
	}

	r := gin.Default()

	modules := []Module{
		NewUserModule(),
	}

	routes.RegisterRoutes(r, getModulRoutes(modules)...)

	return &Application{
		config: cfg,
		router: r,
		modules: modules,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddress)
}

func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
