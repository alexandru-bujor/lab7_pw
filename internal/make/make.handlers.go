package make

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func CreateMakeHandler(c *gin.Context) {
	var make Make
	if err := c.ShouldBindJSON(&make); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := CreateMake(&make); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create make"})
		return
	}
	c.JSON(http.StatusCreated, make)
}

func GetAllMakesHandler(c *gin.Context) {
	makes, err := GetAllMakes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch makes"})
		return
	}
	c.JSON(http.StatusOK, makes)
}

func UpdateMakeHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var make Make
	if err := c.ShouldBindJSON(&make); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := UpdateMake(uint(id), &make); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update make"})
		return
	}
	c.JSON(http.StatusOK, make)
}

func DeleteMakeHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := DeleteMake(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete make"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Make deleted"})
}

func CreateModelHandler(c *gin.Context) {
	var model Model
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := CreateModel(&model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model"})
		return
	}
	c.JSON(http.StatusCreated, model)
}

func UpdateModelHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var model Model
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := UpdateModel(uint(id), &model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model"})
		return
	}
	c.JSON(http.StatusOK, model)
}

func DeleteModelHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := DeleteModel(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete model"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Model deleted"})
}
