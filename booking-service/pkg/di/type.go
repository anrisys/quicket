package di

import (
	"quicket/booking-service/internal/booking"
	"quicket/booking-service/internal/mq/consumer"
	"quicket/booking-service/pkg/config"
)

type App struct {
	Config *config.Config
	Handler *booking.Handler
	EventConsumer *consumer.EventConsumer
}