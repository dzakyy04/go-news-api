package database

import (
	"fmt"
	"go-news-api/models/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed() error {
	if !isDataEmpty() {
		fmt.Println("The database already contains data. Seeder is not executed.")
		return nil
	}

	if err := seedCategories(); err != nil {
		return fmt.Errorf("error seeding categories: %w", err)
	}
	if err := seedUsers(); err != nil {
		return fmt.Errorf("error seeding users: %w", err)
	}
	if err := seedArticles(); err != nil {
		return fmt.Errorf("error seeding articles: %w", err)
	}

	fmt.Println("Successfully seeded the database.")
	return nil
}

func isDataEmpty() bool {
	var count int64

	DB.Model(&entity.User{}).Count(&count)
	if count > 0 {
		return false
	}

	DB.Model(&entity.Category{}).Count(&count)
	if count > 0 {
		return false
	}

	DB.Model(&entity.Article{}).Count(&count)
	return count <= 0
}

func seedCategories() error {
	categories := []entity.Category{
		{Name: "Education", Description: "Related to education, school, and college"},
		{Name: "Entertainment", Description: "Related to entertainment, movies, and series"},
		{Name: "Health", Description: "Related to health, medical, and fitness"},
		{Name: "Music", Description: "Related to music, songs, and albums"},
		{Name: "Technology", Description: "Related to technology, programming, and computing"},
	}

	for _, category := range categories {
		var existingCategory entity.Category
		if err := DB.Where("name = ?", category.Name).First(&existingCategory).Error; err == gorm.ErrRecordNotFound {
			if err := DB.Create(&category).Error; err != nil {
				return fmt.Errorf("error creating category %s: %w", category.Name, err)
			}
			fmt.Printf("Seeded category: %s\n", category.Name)
		}
	}

	return nil
}

func seedUsers() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	users := []entity.User{
		{Name: "Gojo Satoru", Email: "gojo@gmail.com", Password: string(hashedPassword), IsVerified: true},
		{Name: "Ryomen Sukuna", Email: "sukuna@gmailcom", Password: string(hashedPassword), IsVerified: true},
	}

	for _, user := range users {
		var existingUser entity.User
		if err := DB.Where("email = ?", user.Email).First(&existingUser).Error; err == gorm.ErrRecordNotFound {
			if err := DB.Create(&user).Error; err != nil {
				return fmt.Errorf("error creating user %s: %w", user.Email, err)
			}
			fmt.Printf("Seeded user: %s\n", user.Name)
		}
	}

	return nil
}

func seedArticles() error {
	articles := []entity.Article{
		{
			Title:      "The Future of AI in Education",
			Slug:       "future-of-ai-in-education",
			Thumbnail:  "https://example.com/images/ai-education.jpg",
			Content:    "Artificial Intelligence is revolutionizing the education sector...",
			CategoryID: 1,
			AuthorID:   1,
		},
		{
			Title:      "Top 10 Movies of 2024",
			Slug:       "top-10-movies-2024",
			Thumbnail:  "https://example.com/images/movies-2024.jpg",
			Content:    "2024 has been an exceptional year for cinema. Here are our top picks...",
			CategoryID: 2,
			AuthorID:   2,
		},
		{
			Title:      "Breakthrough in Cancer Research",
			Slug:       "breakthrough-cancer-research",
			Thumbnail:  "https://example.com/images/cancer-research.jpg",
			Content:    "Scientists have made a groundbreaking discovery in cancer treatment...",
			CategoryID: 3,
			AuthorID:   1,
		},
		{
			Title:      "The Rise of K-Pop Globally",
			Slug:       "rise-of-kpop-globally",
			Thumbnail:  "https://example.com/images/kpop.jpg",
			Content:    "K-Pop has taken the world by storm. We explore its global impact...",
			CategoryID: 4,
			AuthorID:   2,
		},
		{
			Title:      "Quantum Computing: A New Era",
			Slug:       "quantum-computing-new-era",
			Thumbnail:  "https://example.com/images/quantum-computing.jpg",
			Content:    "Quantum computing is set to revolutionize technology as we know it...",
			CategoryID: 5,
			AuthorID:   1,
		},
	}

	for _, article := range articles {
		var existingArticle entity.Article
		if err := DB.Where("slug = ?", article.Slug).First(&existingArticle).Error; err == gorm.ErrRecordNotFound {
			if err := DB.Create(&article).Error; err != nil {
				return fmt.Errorf("error creating article %s: %w", article.Title, err)
			}
			fmt.Printf("Seeded article: %s\n", article.Title)
		}
	}

	return nil
}
