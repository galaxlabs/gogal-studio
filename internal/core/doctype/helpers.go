package doctype

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

// slugify converts a human-readable name to a lowercase underscore slug.
// e.g. "Module Def" -> "module_def"
func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	value = re.ReplaceAllString(value, "_")
	value = strings.Trim(value, "_")
	if value == "" {
		return "unknown"
	}
	return value
}

// hashBytes returns the hex-encoded SHA-256 hash of raw.
func hashBytes(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
