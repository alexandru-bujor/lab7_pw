package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategoryTemplateHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	tmpl, err := GetCategoryTemplateByCategoryID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": tmpl})
}

func UpsertCategoryTemplateHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var payload CategoryTemplate
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload.CategoryID = id

	if err := UpsertCategoryTemplate(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save template", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template saved", "template": payload})
}

func DeleteCategoryTemplateHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := DeleteCategoryTemplateByCategoryID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete template", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
}
