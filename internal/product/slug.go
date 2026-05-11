package product

import (
	"fmt"
	"regexp"
	"strings"
)

func GenerateSlug(title string, id int) string {
	slug := strings.ToLower(title)

	replacements := map[string]string{
		"ă": "a", "â": "a", "î": "i", "ș": "s", "ț": "t",
		"Ă": "a", "Â": "a", "Î": "i", "Ș": "s", "Ț": "t",
	}
	for old, new := range replacements {
		slug = strings.ReplaceAll(slug, old, new)
	}

	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	slug = fmt.Sprintf("%s-%d", slug, id)

	return slug
}
