package repository

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"be-ac-mobil-ku/domain"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	var db *gorm.DB
	var err error

	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite" // fallback default
	}

	if driver == "mysql" {
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "127.0.0.1"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "3306"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "root"
		}
		password := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "ac_mobil_ku"
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, dbName)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Printf("Failed to connect to MySQL: %v. Falling back to SQLite...", err)
			driver = "sqlite"
		}
	}

	if driver == "postgres" {
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5420" // default/configured
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "ac_mobil_ku"
		}

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			host, user, password, dbName, port)
		
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Printf("Failed to connect to PostgreSQL: %v. Falling back to SQLite...", err)
			driver = "sqlite"
		}
	}

	if driver == "sqlite" {
		dbPath := os.Getenv("DB_SQLITE_PATH")
		if dbPath == "" {
			dbPath = "ac_mobil_ku.db"
		}
		// Ensure directory exists
		dir := filepath.Dir(dbPath)
		if dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
		
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("Failed to initialize SQLite: %v", err)
		}
	}

	log.Println("Database connection established successfully.")

	// Auto Migration
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Bengkel{},
		&domain.Layanan{},
		&domain.Booking{},
		&domain.Rating{},
	)
	if err != nil {
		log.Fatalf("Database AutoMigration failed: %v", err)
	}
	log.Println("Database AutoMigration finished.")

	return db
}
