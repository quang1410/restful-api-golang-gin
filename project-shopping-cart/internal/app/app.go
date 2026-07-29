package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project-shopping-cart/internal/config"
	"project-shopping-cart/internal/routes"
	"project-shopping-cart/internal/validation"
	"syscall"
	"time"

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
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	quit := make(chan os.Signal, 1)
	// SIGINT: Ctrl+C, SIGTERM: termination signal, SIGHUP: terminal closed
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		log.Printf("Server is running at %s\n", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %s\n", err)
		}
	}()

	<-quit
	log.Println("Shutting down signal received, exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s\n", err)
	}

	log.Println("Server exited gracefully")
	return nil
}

func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
