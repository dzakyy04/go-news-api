package database

import (
	"fmt"
	"go-news-api/models/entity"
)

func MigrateDatabase() {
	err := DB.AutoMigrate(&entity.Category{}, &entity.User{}, &entity.OtpCode{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Successfully migrated the database.")
}
