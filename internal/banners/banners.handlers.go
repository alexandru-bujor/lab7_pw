package banners

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Banner struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ImageURL  string    `json:"image_url"`
	LinkURL   string    `json:"link_url"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BannerService struct {
	db *gorm.DB
}

func NewBannerService(db *gorm.DB) *BannerService {
	return &BannerService{db: db}
}

func (s *BannerService) CreateBanner(c *gin.Context) {
	var payload struct {
		ImageURL  string `json:"image_url" binding:"required"`
		LinkURL   string `json:"link_url"`
		IsActive  *bool  `json:"is_active"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	isActive := true
	if payload.IsActive != nil {
		isActive = *payload.IsActive
	}

	banner := Banner{
		ImageURL:  payload.ImageURL,
		LinkURL:   payload.LinkURL,
		IsActive:  isActive,
		SortOrder: payload.SortOrder,
	}

	if err := s.db.Create(&banner).Error; err != nil {
		log.Printf("Error creating banner: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create banner"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"banner": banner})
}

func (s *BannerService) GetBanners(c *gin.Context) {
	var banners []Banner

	query := s.db.Order("sort_order ASC, id DESC")

	activeOnly := c.Query("active")
	if activeOnly == "1" || activeOnly == "true" {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&banners).Error; err != nil {
		log.Printf("Error fetching banners: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch banners"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"banners": banners})
}

func (s *BannerService) UpdateBanner(c *gin.Context) {
	id := c.Param("id")

	var banner Banner
	if err := s.db.First(&banner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
		return
	}

	var payload struct {
		ImageURL  string `json:"image_url"`
		LinkURL   string `json:"link_url"`
		IsActive  *bool  `json:"is_active"`
		SortOrder *int   `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if payload.ImageURL != "" {
		banner.ImageURL = payload.ImageURL
	}
	if payload.LinkURL != "" {
		banner.LinkURL = payload.LinkURL
	}
	if payload.IsActive != nil {
		banner.IsActive = *payload.IsActive
	}
	if payload.SortOrder != nil {
		banner.SortOrder = *payload.SortOrder
	}

	if err := s.db.Save(&banner).Error; err != nil {
		log.Printf("Error updating banner: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update banner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"banner": banner})
}

func (s *BannerService) DeleteBanner(c *gin.Context) {
	id := c.Param("id")

	if err := s.db.Delete(&Banner{}, id).Error; err != nil {
		log.Printf("Error deleting banner: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete banner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Banner deleted successfully"})
}

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	svc := NewBannerService(db)
	bannersGroup := r.Group("/banners")
	{
		bannersGroup.GET("", svc.GetBanners)
		bannersGroup.GET("/", svc.GetBanners)
		bannersGroup.POST("", svc.CreateBanner)
		bannersGroup.POST("/", svc.CreateBanner)
		bannersGroup.PUT("/:id", svc.UpdateBanner)
		bannersGroup.DELETE("/:id", svc.DeleteBanner)
	}
}
