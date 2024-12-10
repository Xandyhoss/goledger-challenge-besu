package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Contract struct {
	gorm.Model
	ContractNumber uint `gorm:"not null"`
}

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database: %w", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func StartContractDB() {

	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	if (db.Migrator().HasTable(&Contract{})) {
		log.Println("Contract table already exists")
	} else {
		db.AutoMigrate(&Contract{})
		log.Println("Contract table created")
	}

	var contract Contract
	result := db.First(&contract, 1)
	if result.Error == gorm.ErrRecordNotFound {
		contract.ContractNumber = 0
		db.Create(&contract)
		log.Println("Contract with ID 1 created with ContractNumber 0")
	} else {
		log.Println("Contract with ID 1 already exists, skipping creation")
	}
}
