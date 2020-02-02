package cache

import "encoding/base64"

// GenerateName returns base64 encoded string
func GenerateName(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
