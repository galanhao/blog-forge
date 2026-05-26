// Package util provides shared utility functions.
package util

import "unicode"

// Slugify converts a string to a URL-friendly slug.
// It lowercases ASCII letters, converts spaces/underscores to hyphens,
// keeps digits and CJK characters, and drops everything else.
func Slugify(s string) string {
	result := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			result = append(result, r)
		case r >= 'A' && r <= 'Z':
			result = append(result, r+32) // to lower
		case r == ' ' || r == '_':
			result = append(result, '-')
		case unicode.Is(unicode.Han, r):
			result = append(result, r)
		}
	}
	return string(result)
}
