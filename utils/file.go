package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

var AllowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

func IsImage(file *multipart.FileHeader) bool {
	extension := filepath.Ext(file.Filename)
	_, ok := AllowedImageTypes[extension]
	return ok
}

func SaveImageFile(ctx *fiber.Ctx, key string, path string) (string, error) {
	file, err := ctx.FormFile(key)
	if err != nil {
		return "", err
	}

	// Validate file type
	if !IsImage(file) {
		return "", errors.New("invalid file type, only images are allowed")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}

	// Generate file path
	fileName := fmt.Sprintf("%d%s", time.Now().Unix(), filepath.Ext(file.Filename))
	filePath := filepath.ToSlash(filepath.Join(path, fileName))

	// Save file
	if err := ctx.SaveFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}
