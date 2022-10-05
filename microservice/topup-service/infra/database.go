package infra

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitializeDatabase() {
	fmt.Println("Start Connecting to DB...")

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable ",
		os.Getenv("db_host"),
		os.Getenv("db_username"),
		os.Getenv("db_password"),
		os.Getenv("db_name"),
		os.Getenv("db_port"),
	)

	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		fmt.Println("Failed to open connection")
		panic(err)
	}

	fmt.Println("Success to connect DB")

	DB = db
}
