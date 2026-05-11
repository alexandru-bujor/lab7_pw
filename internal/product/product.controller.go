package product

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	products, err := GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
		return
	}

	var publicProducts []Product
	for _, p := range products {
		if p.AccountingType == "" || p.AccountingType == "on_book" {
			publicProducts = append(publicProducts, p)
		}
	}

	if q := strings.TrimSpace(c.Query("q")); q != "" {
		lower := strings.ToLower(q)
		filtered := make([]Product, 0, len(publicProducts))
		for _, p := range publicProducts {
			if strings.Contains(strings.ToLower(p.Title), lower) ||
				strings.Contains(strings.ToLower(p.Slug), lower) ||
				strings.Contains(strings.ToLower(p.BrandName), lower) {
				filtered = append(filtered, p)
			}
		}
		publicProducts = filtered
	}

	c.JSON(http.StatusOK, gin.H{
		"products": publicProducts,
	})
}

func CreateProduct(c *gin.Context) {
	var input Product

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := AddProduct(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product created successfully",
		"product": input,
	})
}

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

	var publicVariants []Product
	for _, v := range variants {
		if v.AccountingType == "" || v.AccountingType == "on_book" {
			publicVariants = append(publicVariants, v)
		}
	}

	c.JSON(http.StatusOK, gin.H{"variants": publicVariants})
}

func GetProductBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")
	p, err := GetProductBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": p})
}
