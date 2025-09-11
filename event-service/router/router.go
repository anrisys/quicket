package router

import (
	"fmt"
	"time"

	"github.com/anrisys/quicket/event-service/internal"
	"github.com/anrisys/quicket/event-service/pkg/di"
	"github.com/anrisys/quicket/event-service/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(app *di.EventServiceApp) *gin.Engine {
	r := gin.New()

	// Logging
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		internal.RegisterCustomValidation(v)
	}

	// Middlewares
	r.Use(middleware.ZerologLogger(), gin.Recovery(), middleware.ErrorMiddleware())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	registerRoutes(r, app)

	return r
}

func registerRoutes(r *gin.Engine, app *di.EventServiceApp) {
	public := r.Group("/api/v1/events")
	public.GET("/:publicID", app.Handler.GetEventByPublicID)
	
	protected := r.Group("/api/v1/events")
	protected.Use(middleware.JWTAuthMiddleware(app.Config.Security.JWTSecret))
	{
		protected.POST("/", app.Handler.Create)
	}
}