package router

import (
	"quicket/booking-service/internal"
	"quicket/booking-service/internal/booking"
	"quicket/booking-service/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(h *booking.Handler, l zerolog.Logger, jwtSecret string) *gin.Engine {
	r := gin.New()

	// Custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		internal.RegisterCustomValidation(v)
	}

	// Middlewares
	r.Use(
		middleware.HTTPLogger(l),
		gin.Recovery(),
		middleware.ErrorMiddleware(),
	)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	r.GET("/api/v1/bookings/health", h.HealthCheck)
	protected := r.Group("/api/v1/bookings")
	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		protected.POST("/", h.CreateBooking)
	}

	// TODO: move to payment-service once extracted
	// Endpoint path intentionally uses /payments prefix to reflect domain ownership
	r.POST("/api/v1/payment/webhook", h.ReceivePaymentWebhook)

	admin := r.Group("/api/v1/admin/bookings")
	admin.Use(middleware.JWTAuthMiddleware(jwtSecret))

	return r
}
