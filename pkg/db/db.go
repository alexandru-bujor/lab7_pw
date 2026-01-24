package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// Validate required environment variables
	if user == "" {
		log.Fatal("❌ DB_USER environment variable is not set")
	}
	// Note: DB_PASS can be empty if MySQL user has no password
	if host == "" {
		log.Fatal("❌ DB_HOST environment variable is not set")
	}
	if port == "" {
		port = "3306" // Default MySQL port
	}
	if name == "" {
		log.Fatal("❌ DB_NAME environment variable is not set")
	}

	passwordInfo := "with password"
	if pass == "" {
		passwordInfo = "without password"
	}
	log.Printf("🔄 Attempting to connect to MySQL at %s:%s (database: %s, user: %s, %s)", host, port, name, user, passwordInfo)

	// Build DSN - handle empty password correctly
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to connect to MySQL database")
		log.Printf("   Host: %s:%s", host, port)
		log.Printf("   Database: %s", name)
		log.Printf("   User: %s", user)
		log.Printf("   Error: %v", err)
		log.Println("\n💡 Troubleshooting tips:")
		log.Println("   1. Verify MySQL is running: mysql -u root -p")
		log.Println("   2. Check if the database exists: CREATE DATABASE IF NOT EXISTS megamobile;")
		log.Println("   3. Verify credentials in .env file match your MySQL setup")
		log.Println("   4. If using Docker: docker-compose up -d mysql")
		log.Fatal("")
	}

	DB = connection
	log.Println("📡 MySQL connected successfully")
}

// HealthCheck tries to ping the DB connection.
// Returns true if DB is reachable, false otherwise.
func HealthCheck() bool {
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		return false
	}

	return true
}
