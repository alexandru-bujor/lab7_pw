package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /products
func GetProducts(c *gin.Context) {
	products, err := GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
		return
	}

	// Filter out off-book products for public API (clients should not see them)
	var publicProducts []Product
	for _, p := range products {
		// Only include products that are on-book or don't have accounting_type set
		if p.AccountingType == "" || p.AccountingType == "on_book" {
			publicProducts = append(publicProducts, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": publicProducts,
	})
}

// POST /products
func CreateProduct(c *gin.Context) {
	var input Product

	// Validate + bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create product + nested images + nested variants
	if err := AddProduct(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product created successfully",
		"product": input,
	})
}

// DELETE /products/:id
func DeleteProductHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// PUT /products/:id
func UpdateProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := UpdateProduct(id, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": input,
	})
}

// GET /products/:id/variants
func GetProductVariantsHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	variants, err := GetProductVariants(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch variants"})
		return
	}

	// Filter out off-book variants for public API
	var publicVariants []Product
	for _, v := range variants {
		if v.AccountingType == "" || v.AccountingType == "on_book" {
			publicVariants = append(publicVariants, v)
		}
	}

	c.JSON(http.StatusOK, gin.H{"variants": publicVariants})
}
