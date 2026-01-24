package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategoriesHandler(c *gin.Context) {
	categories, err := GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func CreateCategoryHandler(c *gin.Context) {
	var payload Category

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := AddCategory(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "category created",
		"category": payload,
	})
}

func UpdateCategoryHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var payload Category
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload.ID = id

	if err := UpdateCategory(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "category updated",
		"category": payload,
	})
}

func DeleteCategoryHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}
