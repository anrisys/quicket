package main

import (
	"context"
	"fmt"
	"log"
	"quicket/booking-service/pkg/di"
	"quicket/booking-service/router"
	"runtime/debug"
)

// @title Quicket Bookings Service API
// @version 1.0
// @description Bookings service API

// @contact.name Quicket Support
// @contact.url https://github.com/anrisys/quicket
// @contact.email anris.y.simorangkir@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

// @host localhost:8091
// @BasePath /api/v1/bookings
func main() {
	app, err := di.InitializeApp()
    if err != nil {
        log.Fatal(err)
        log.Printf("Full stack trace:\n%s", debug.Stack())
        log.Fatalf("Failed to initialize app: %v", err)
    }

    go func ()  {
        if err := app.EventConsumer.Start(context.Background()); err != nil {
            log.Fatalf("Failed to start event consumer: %v", err)
        }
    }()

    go func ()  {
        if err := app.UserConsumer.Start(context.Background()); err != nil {
            log.Fatalf("Failed to start event consumer: %v", err)
        }
    }()
    
    r := router.SetupRouter(app)
    
    addr := fmt.Sprintf(":%s", app.Config.Server.Port)
	if err := r.Run(addr); err != nil {
        log.Fatalf("Server failed :%v", err)
    }
}