package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifySecondaryPasswordHandler verifies the secondary panel password
func VerifySecondaryPasswordHandler(c *gin.Context) {
	var body struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	valid, err := VerifySecondaryPanelPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify password"})
		return
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// ChangeSecondaryPasswordHandler changes the secondary panel password
func ChangeSecondaryPasswordHandler(c *gin.Context) {
	var body struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=4"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify current password
	valid, err := VerifySecondaryPanelPassword(body.CurrentPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify current password"})
		return
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid current password"})
		return
	}

	// Set new password
	if err := SetSecondaryPanelPassword(body.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
