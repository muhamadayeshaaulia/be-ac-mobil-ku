package routes

import (
	"net/http"
	"os"

	"be-ac-mobil-ku/delivery"
	"be-ac-mobil-ku/domain"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App            *gin.Engine
	UserUC         domain.UserUsecase
	BengkelUC      domain.BengkelUsecase
	LayananUC      domain.LayananUsecase
	BookingUC      domain.BookingUsecase
	RatingUC       domain.RatingUsecase
	RecUC          domain.RecommendationUsecase
	AuthMiddleware gin.HandlerFunc
	DBName         string
}

func SetupRoutes(cfg RouteConfig) {
	// Inisialisasi Handlers
	userHandler := delivery.NewUserHandler(cfg.UserUC)
	bengkelHandler := delivery.NewBengkelHandler(cfg.BengkelUC)
	layananHandler := delivery.NewLayananHandler(cfg.LayananUC, cfg.BengkelUC)
	bookingHandler := delivery.NewBookingHandler(cfg.BookingUC)
	ratingHandler := delivery.NewRatingHandler(cfg.RatingUC)
	recHandler := delivery.NewRecommendationHandler(cfg.RecUC)

	// 1. Public Endpoint (Tanpa Autentikasi)
	cfg.App.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"message":     "Backend AC Mobil Ku siap digunakan!",
			"environment": os.Getenv("APP_ENV"),
			"database":    cfg.DBName,
		})
	})

	// 2. Protected Endpoints (Perlu JWT Firebase Auth Bearer Token)
	api := cfg.App.Group("/api")
	api.Use(cfg.AuthMiddleware)
	{
		// --- USER ROUTES ---
		api.GET("/user/profile", userHandler.GetProfile)
		api.POST("/user/profile", userHandler.RegisterOrUpdate)
		api.POST("/user/location", userHandler.UpdateLocation)

		// --- BENGKEL ROUTES ---
		api.GET("/bengkel", bengkelHandler.ListBengkel)
		api.GET("/bengkel/detail/:id", bengkelHandler.GetBengkelByID)
		api.GET("/bengkel/my", bengkelHandler.GetMyBengkel)
		api.POST("/bengkel", bengkelHandler.CreateBengkel)
		api.PUT("/bengkel", bengkelHandler.UpdateBengkel)

		// --- LAYANAN CATALOG ROUTES ---
		api.GET("/layanan/bengkel/:bengkel_id", layananHandler.GetLayananByBengkel)
		api.GET("/layanan/detail/:id", layananHandler.GetLayananByID)
		api.POST("/layanan", layananHandler.CreateLayanan)
		api.PUT("/layanan/:id", layananHandler.UpdateLayanan)
		api.DELETE("/layanan/:id", layananHandler.DeleteLayanan)

		// --- BOOKING SERVICE ROUTES ---
		api.POST("/booking", bookingHandler.CreateBooking)
		api.GET("/booking/history", bookingHandler.GetHistory)
		api.GET("/booking/queue", bookingHandler.GetQueue)
		api.GET("/booking/slots", bookingHandler.GetFullTimeSlots)
		api.PUT("/booking/:id/status", bookingHandler.UpdateStatus)

		// --- RATING & REVIEW ROUTES ---
		api.POST("/rating", ratingHandler.AddRating)
		api.GET("/rating/bengkel/:bengkel_id", ratingHandler.GetRatings)

		// --- REKOMENDASI HYBRID ROUTES ---
		api.GET("/recommendations", recHandler.GetRecommendations)
	}
}
