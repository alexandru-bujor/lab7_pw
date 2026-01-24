package product

import (
	"regexp"
	"strings"
	"fmt"
)

// GenerateSlug creates a URL-friendly slug from a title
func GenerateSlug(title string, id int) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace Romanian characters
	replacements := map[string]string{
		"ă": "a", "â": "a", "î": "i", "ș": "s", "ț": "t",
		"Ă": "a", "Â": "a", "Î": "i", "Ș": "s", "Ț": "t",
	}
	for old, new := range replacements {
		slug = strings.ReplaceAll(slug, old, new)
	}
	
	// Remove special characters and replace spaces with hyphens
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	
	// Always append ID for uniqueness (important for variants with same title)
	slug = fmt.Sprintf("%s-%d", slug, id)
	
	return slug
}
