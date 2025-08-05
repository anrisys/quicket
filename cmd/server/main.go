package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anrisys/quicket/pkg/di"
	"github.com/anrisys/quicket/pkg/middleware"
	"github.com/gin-gonic/gin"
)

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
	
	router.Use(middleware.ZerologLogger(), gin.Recovery(), middleware.ErrorHandler())
    registerRoutes(router, app)
    
    addr := fmt.Sprintf(":%s", app.Config.Server.Port)

	if err := router.Run(addr); err != nil {
        log.Fatalf("Server failed :%v", err)
    }
}

func registerRoutes(r *gin.Engine, app *di.App)  {
    r.POST("/register", app.UserHandler.Register)
    r.POST("/login", app.UserHandler.Login)
}