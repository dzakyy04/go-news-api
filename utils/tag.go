package utils

import (
	"go-news-api/database"
	"go-news-api/models/entity"
)

func CreateOrFindTags(tagNames []string) ([]entity.Tag, error) {
	var tags []entity.Tag

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, tagName := range tagNames {
		var tag entity.Tag
		if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, entity.Tag{Name: tagName}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func AssociateTagsWithArticle(articleID uint, tags []entity.Tag) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var article entity.Article
	if err := tx.First(&article, articleID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&article).Association("Tags").Replace(tags); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
