package database

import (
	"fmt"
	"kyaho-space/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection, accessible by all packages
var DB *gorm.DB

// OpenConnection initializes the global DB connection using the provided config
func OpenConnection(cfg config.Config) *gorm.DB {

	dsn := cfg.DSN()
	fmt.Println("Connecting to database...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Database connection error:", err)
		panic(err)
	}

	fmt.Println("Database connected successfully!")
	DB = db
	return db
}

// MigrateSchema auto-migrates all entity models
// Add new models here as the project grows (e.g. &entities.User{}, &entities.Application{})
func MigrateSchema(models ...interface{}) {
	if DB == nil {
		panic("Database not initialized. Call OpenConnection() first.")
	}

	err := DB.AutoMigrate(models...)
	if err != nil {
		panic(fmt.Sprintf("Migration failed: %v", err))
	}

	fmt.Println("Database migration completed!")
}
