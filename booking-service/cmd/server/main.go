package main

import (
	"fmt"
	"quicket/booking-service/internal/booking"
	"quicket/booking-service/internal/infrastructure"
	"quicket/booking-service/internal/logger"
	"quicket/booking-service/internal/snapshot"
	"quicket/booking-service/router"

	amqp "github.com/rabbitmq/amqp091-go"
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
	logger := logger.NewLogger()

	cfg, err := infrastructure.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed load configuration")
	}

	// DB
	db, err := infrastructure.NewMySQL(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize database")
	}

	// RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQ.URL())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect rabbitmq")
	}

	rabbitmqPublisher, err := infrastructure.NewRabbitMQPublisher(
		conn,
		cfg.RabbitMQ.BookingServiceExchange,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize rabbitmq publisher")
	}

	// Repositories
	bwr := booking.NewBookingWriteRepoImpl(db, logger)
	brr := booking.NewBookingReadRepo(db, logger)

	// Clients
	ecs := infrastructure.NewHTTPEventClient(cfg.Clients.EventServiceURL)

	// Snapshot
	esr := snapshot.NewEvSnapshotRepoImpl(db, logger)
	ess := snapshot.NewEventSnapshotService(esr, logger)

	usr := snapshot.NewUserSnapshotRepoImpl(db, logger)
	uss := snapshot.NewUserSnapshotService(usr, logger)

	// Usecases
	usu := booking.NewUserUsecase(
		bwr,
		brr,
		logger,
		ecs,
		rabbitmqPublisher,
		ess,
		uss,
	)

	adu := booking.NewAdminUsecase(brr, logger)

	// Handler
	handler := booking.NewHandler(usu, adu)

	// Router
	r := router.SetupRouter(handler, logger, cfg.JWT.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	logger.Info().
		Str("addr", addr).
		Msg("starting server")

	if err := r.Run(addr); err != nil {
		logger.Fatal().Err(err).Msg("server failed")
	}
}
