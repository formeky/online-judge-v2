package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"online-judge/internal/api/handler"
	"online-judge/internal/api/middleware"
)

func SetupRouter(
	logger *zap.Logger,
	problemHandler *handler.ProblemHandler,
	submissionHandler *handler.SubmissionHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(logger))

	r.GET("/health", handler.Health)

	v1 := r.Group("/api/v1")
	{
		problems := v1.Group("/problems")
		{
			problems.GET("", problemHandler.List)
			problems.GET("/:id", problemHandler.Get)

			admin := problems.Group("", middleware.RequireAuth(), middleware.RequireAdmin())
			{
				admin.POST("", problemHandler.Create)
				admin.PUT("/:id", problemHandler.Update)
				admin.DELETE("/:id", problemHandler.Delete)
				admin.POST("/:id/testcases", problemHandler.AddTestCase)
			}
		}

		submissions := v1.Group("/submissions", middleware.RequireAuth())
		{
			submissions.POST("", submissionHandler.Submit)
			submissions.GET("", submissionHandler.List)
			submissions.GET("/:id", submissionHandler.Get)
		}
	}

	return r
}
