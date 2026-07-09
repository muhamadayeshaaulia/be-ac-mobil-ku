package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"be-ac-mobil-ku/middleware"
	"be-ac-mobil-ku/repository"
	"be-ac-mobil-ku/routes"
	"be-ac-mobil-ku/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables.")
	}

	// 1. Initialize environment configurations (Dev mode default for database / auth bypass)
	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "development")
	}

	log.Printf("Starting AC Mobil Ku Backend in [%s] mode...", os.Getenv("APP_ENV"))

	// 2. Initialize database connection & schemas
	db := repository.InitDB()

	// 3. Initialize Firebase Auth middleware
	authMiddleware := middleware.InitAuthMiddleware()
	verifyTokenFunc := authMiddleware.VerifyToken()

	// 4. Initialize Repository Layer
	userRepo := repository.NewGormUserRepository(db)
	bengkelRepo := repository.NewGormBengkelRepository(db)
	layananRepo := repository.NewGormLayananRepository(db)
	bookingRepo := repository.NewGormBookingRepository(db)
	ratingRepo := repository.NewGormRatingRepository(db)

	// 5. Initialize Usecase Layer
	userUC := usecase.NewUserUsecase(userRepo)
	bengkelUC := usecase.NewBengkelUsecase(bengkelRepo, ratingRepo)
	layananUC := usecase.NewLayananUsecase(layananRepo)
	bookingUC := usecase.NewBookingUsecase(bookingRepo, bengkelRepo)
	ratingUC := usecase.NewRatingUsecase(ratingRepo, bookingRepo)
	recUC := usecase.NewRecommendationUsecase(userRepo, bengkelRepo, ratingRepo, layananRepo)

	// 6. Initialize Gin Web Framework
	r := gin.Default()

	// Enable CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Setup static assets serving for uploads folder
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		_ = os.Mkdir("uploads", 0755)
	}
	r.Static("/uploads", "./uploads")

	// REST API endpoint for uploading files locally
	r.POST("/api/upload", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
		filepath := fmt.Sprintf("uploads/%s", filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		host := c.Request.Host
		fileURL := fmt.Sprintf("http://%s/uploads/%s", host, filename)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"url":    fileURL,
		})
	})

	// 8. Register Handlers (Delivery Layer)
	routes.SetupRoutes(routes.RouteConfig{
		App:            r,
		UserUC:         userUC,
		BengkelUC:      bengkelUC,
		LayananUC:      layananUC,
		BookingUC:      bookingUC,
		RatingUC:       ratingUC,
		RecUC:          recUC,
		AuthMiddleware: verifyTokenFunc,
		DBName:         db.Name(),
	})

	// 9. Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening and serving on port :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}
}