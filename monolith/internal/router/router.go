package router

import (
	"fmt"
	"time"

	"github.com/anrisys/quicket/internal/validation"
	"github.com/anrisys/quicket/pkg/di"
	"github.com/anrisys/quicket/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(app *di.App) *gin.Engine {
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
		validation.RegisterCustomValidation(v)
	}

	// Middlewares
	r.Use(middleware.ZerologLogger(), gin.Recovery(), middleware.ErrorHandler())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	registerRoutes(r, app)

	return r
}

func registerRoutes(r *gin.Engine, app *di.App) {
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuthMiddleware(app.Config.Security.JWTSecret))
	{
		events := protected.Group("/events")
		events.Use(middleware.AuthorizedRole([]string{"admin", "organizer"}))
		{
			events.POST("", app.EventHandler.Create)
		}

		bookings := protected.Group("/bookings")
		{
			bookings.POST(":eventID", app.BookingHandler.Create)
		}
	}
}
