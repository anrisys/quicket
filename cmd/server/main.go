package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anrisys/quicket/internal/validation"
	"github.com/anrisys/quicket/pkg/di"
	"github.com/anrisys/quicket/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title Quicket API
// @version 1.0
// @description Event Booking and Management System API

// @contact.name Quicket Support
// @contact.url https://github.com/anrisys/quicket
// @contact.email your.email@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

// @host localhost:8080
// @BasePath /api/v1
func main() {
	app, err := di.InitializeApp()
    if err != nil {
        log.Fatalf("Failed to initialize app: %v", err)
    }
    
    router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
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

    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        validation.RegisterCustomValidation(v)
    }
	
	router.Use(middleware.ZerologLogger(), gin.Recovery(), middleware.ErrorHandler())

    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    registerRoutes(router, app)
    
    addr := fmt.Sprintf(":%s", app.Config.Server.Port)

	if err := router.Run(addr); err != nil {
        log.Fatalf("Server failed :%v", err)
    }
}

func registerRoutes(r *gin.Engine, app *di.App)  {
    // r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    public := r.Group("/api/v1")
    {
        public.POST("/register", app.UserHandler.Register)
        public.POST("/login", app.UserHandler.Login)
    }

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