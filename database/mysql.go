package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// getEnv environment variable ile fallback (MySQL için)
func getEnvMySQL(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// ConnectMysql was updated for Docker environment
func ConnectMysql() {
	// Environment variable'dan MySQL config al
	mysqlHost := getEnvMySQL("MYSQL_HOST", "localhost")
	mysqlPort := getEnvMySQL("MYSQL_PORT", "3306")
	mysqlDatabase := getEnvMySQL("MYSQL_DATABASE", "gps_tracker")
	mysqlUser := getEnvMySQL("MYSQL_USER", "root")
	mysqlPassword := getEnvMySQL("MYSQL_PASSWORD", "password")

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)

	var err error
	// Connect with retry logic (MySQL başlaması zaman alabilir)
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err == nil {
			// Connection'ı test et
			sqlDB, err := DB.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					break
				}
			}
		}

		log.Printf("MySQL connection attempt %d failed: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * 2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to MySQL after retries:", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("Connected to MySQL successfully at: %s:%s, Database: %s", mysqlHost, mysqlPort, mysqlDatabase)
}
