package main

import (
	"fmt"
	"log"

	"github.com/anrisys/quicket/event-service/pkg/di"
	"github.com/anrisys/quicket/event-service/router"
)

// @title Quicket Event Service API
// @version 1.0
// @description Event service API

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
// @BasePath /api/v1
func main() {
	app, err := di.InitializeEventServiceApp()
    if err != nil {
        log.Fatalf("Failed to initialize app: %v", err)
    }
    
    r := router.SetupRouter(app)
    
    addr := fmt.Sprintf(":%s", app.Config.Server.Port)
	if err := r.Run(addr); err != nil {
        log.Fatalf("Server failed :%v", err)
    }
}