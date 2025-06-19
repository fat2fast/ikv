package shared

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb() *gorm.DB {
	log.Println("Db Connecting .... ")
	host := os.Getenv("DATABASE_HOST")
	dbName := os.Getenv("DATABASE_DBNAME")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	schema := os.Getenv("DATABASE_SCHEMA")
	timezone := os.Getenv("TIMEZONE")

	// Xây dựng connection string cho GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, dbName, port, timezone)

	// Thiết lập POSTGRESQL_URL cho migration
	migrateDsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table-quoted=%s",
		user, password, host, port, dbName, schema)

	if err := os.Setenv("POSTGRESQL_URL", migrateDsn); err != nil {
		log.Printf("Error setting POSTGRESQL_URL: %v", err)
	}

	// Kết nối database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect database:", err)
	}

	// Thiết lập schema
	db.Exec("SET search_path TO " + schema)
	log.Printf("Connected to database with schema: %s", schema)

	return db
}
func InitDb(db *gorm.DB) gin.HandlerFunc {
	log.Println("Init Db to Context.... ")
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
