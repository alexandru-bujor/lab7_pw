
package servicecategory

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func UpdateServiceCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	var req ServiceCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := UpdateServiceCategory(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}
	c.JSON(http.StatusOK, req)
}

func DeleteServiceCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	if err := DeleteServiceCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

type createCategoryRequest struct {
	Name struct {
		RO string `json:"ro"`
		RU string `json:"ru"`
		EN string `json:"en"`
	} `json:"name"`
}

func CreateServiceCategoryHandler(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sc := ServiceCategory{
		Name: req.Name,
	}
	if err := CreateServiceCategory(&sc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, sc)
}

func GetAllServiceCategoriesHandler(c *gin.Context) {
	categories, err := GetAllServiceCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}
