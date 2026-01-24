package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetServices(c *gin.Context) {
	services, err := GetAllServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}
	c.JSON(http.StatusOK, services)
}

func GetServicesByPartHandler(c *gin.Context) {
	partStr := c.Query("part")
	part, err := strconv.Atoi(partStr)
	if err != nil || (part != 1 && part != 2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service part (must be 1 or 2)"})
		return
	}

	services, err := GetServicesByPart(part)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}
	c.JSON(http.StatusOK, services)
}

func CreateServiceHandler(c *gin.Context) {
	var service Service
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if service.ServicePart != 1 && service.ServicePart != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_part must be 1 or 2"})
		return
	}

	if err := CreateService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}
	c.JSON(http.StatusCreated, service)
}

func UpdateServiceHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var service Service
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := UpdateService(uint(id), &service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}
	c.JSON(http.StatusOK, service)
}

func DeleteServiceHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := DeleteService(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted"})
}
