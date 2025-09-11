package main

import (
	"fmt"
	"log"
	"quicket/booking-service/pkg/di"
	"quicket/booking-service/router"
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

// @host localhost:8082
// @BasePath /api/v1/bookings
func main() {
	app, err := di.InitializeApp()
    if err != nil {
        log.Fatalf("Failed to initialize app: %v", err)
    }
    
    r := router.SetupRouter(app)
    
    addr := fmt.Sprintf(":%s", app.Config.Server.Port)
	if err := r.Run(addr); err != nil {
        log.Fatalf("Server failed :%v", err)
    }
}