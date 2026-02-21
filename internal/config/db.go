package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(env Env) {
	if env.DBPass == "" {
		log.Fatal("DB_PASSWORD missing in .env")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Dhaka",
		env.DBHost, env.DBUser, env.DBPass, env.DBName, env.DBPort, env.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect DB: ", err)
	}

	// verify DB connection is really working
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB: ", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("DB ping failed: ", err)
	}

	DB = db
	log.Println("DB connected")
}