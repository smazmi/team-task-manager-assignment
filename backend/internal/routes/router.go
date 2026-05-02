package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/handlers"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/middleware"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	projectHandler *handlers.ProjectHandler,
	taskHandler *handlers.TaskHandler,
	dashboardHandler *handlers.DashboardHandler,
	authMiddleware *middleware.AuthMiddleware,
	rbacMiddleware *middleware.RBACMiddleware,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	api := router.Group("/api")
	{
		authGroup := api.Group("/auth")
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())

		projectGroup := protected.Group("/projects")
		projectGroup.GET("", projectHandler.ListProjects)
		projectGroup.POST("", projectHandler.CreateProject)
		projectGroup.GET("/:projectId", rbacMiddleware.RequireProjectMember(), projectHandler.GetProject)
		projectGroup.PATCH("/:projectId", rbacMiddleware.RequireProjectAdmin(), projectHandler.UpdateProject)
		projectGroup.DELETE("/:projectId", rbacMiddleware.RequireProjectAdmin(), projectHandler.DeleteProject)
		projectGroup.GET("/:projectId/members", rbacMiddleware.RequireProjectMember(), projectHandler.ListMembers)
		projectGroup.POST("/:projectId/members", rbacMiddleware.RequireProjectAdmin(), projectHandler.AddMember)
		projectGroup.DELETE("/:projectId/members/:userId", rbacMiddleware.RequireProjectAdmin(), projectHandler.RemoveMember)

		taskGroup := protected.Group("/tasks")
		taskGroup.POST("", rbacMiddleware.RequireProjectAdmin(), taskHandler.CreateTask)
		taskGroup.GET("", rbacMiddleware.RequireProjectMember(), taskHandler.ListTasks)
		taskGroup.GET("/:taskId", rbacMiddleware.RequireTaskActor(), taskHandler.GetTask)
		taskGroup.PATCH("/:taskId", rbacMiddleware.RequireTaskAdmin(), taskHandler.UpdateTask)
		taskGroup.DELETE("/:taskId", rbacMiddleware.RequireTaskAdmin(), taskHandler.DeleteTask)
		taskGroup.PATCH("/:taskId/status", rbacMiddleware.RequireTaskActor(), taskHandler.UpdateTaskStatus)

		dashboardGroup := protected.Group("/dashboard")
		dashboardGroup.GET("/stats", rbacMiddleware.RequireProjectMember(), dashboardHandler.GetStats)
	}

	return router
}
