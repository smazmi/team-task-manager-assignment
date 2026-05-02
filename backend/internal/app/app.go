package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/config"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/db"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/handlers"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/middleware"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/repository"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/routes"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/service"
	"gorm.io/gorm"
)

type App struct {
	config   config.Config
	database *gorm.DB
	server   *http.Server
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	gin.SetMode(cfg.GinMode)

	database, err := db.NewPostgres(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if cfg.AutoMigrate {
		if err := db.AutoMigrate(database); err != nil {
			_ = closeDatabase(database)
			return nil, fmt.Errorf("auto-migrate database: %w", err)
		}
	}

	router := buildRouter(cfg, database)
	server := newHTTPServer(cfg.AppPort, router)

	return &App{
		config:   cfg,
		database: database,
		server:   server,
	}, nil
}

func Run() error {
	application, err := New()
	if err != nil {
		return err
	}

	return application.Run()
}

func (a *App) Run() error {
	defer func() {
		if err := closeDatabase(a.database); err != nil {
			log.Printf("database close failed: %v", err)
		}
	}()

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server listening on :%s", a.config.AppPort)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
		close(serverErrors)
	}()

	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err, ok := <-serverErrors:
		if ok && err != nil {
			return fmt.Errorf("server failed: %w", err)
		}
		return nil
	case <-signalContext.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	return nil
}

func buildRouter(cfg config.Config, database *gorm.DB) *gin.Engine {
	userRepository := repository.NewUserRepository(database)
	projectRepository := repository.NewProjectRepository(database)
	taskRepository := repository.NewTaskRepository(database)

	authService := service.NewAuthService(userRepository, cfg.JWTSecret, cfg.JWTTTL)
	projectService := service.NewProjectService(projectRepository, userRepository)
	taskService := service.NewTaskService(taskRepository, projectRepository)
	dashboardService := service.NewDashboardService(projectRepository, taskRepository)

	authHandler := handlers.NewAuthHandler(authService)
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	rbacMiddleware := middleware.NewRBACMiddleware(projectRepository, taskRepository)

	return routes.SetupRouter(
		authHandler,
		projectHandler,
		taskHandler,
		dashboardHandler,
		authMiddleware,
		rbacMiddleware,
	)
}

func newHTTPServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func closeDatabase(database *gorm.DB) error {
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
