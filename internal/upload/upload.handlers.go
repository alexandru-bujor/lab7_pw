package upload

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadProductImageHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, webp, gif"})
		return
	}

	timestamp := time.Now().Unix()
	uniqueID := uuid.New().String()[:8]
	filename := fmt.Sprintf("%d_%s%s", timestamp, uniqueID, ext)
	savePath := filepath.Join("uploads", "products", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image", "details": err.Error()})
		return
	}

	imageURL := fmt.Sprintf("/uploads/products/%s", filename)
	if gin.Mode() == gin.DebugMode {
		scheme := "http"
		if c.GetHeader("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		host := c.GetHeader("Host")
		if host == "" {
			host = c.Request.Host
		}
		if host != "" {
			imageURL = fmt.Sprintf("%s://%s%s", scheme, host, imageURL)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"url":     imageURL,
	})
}
