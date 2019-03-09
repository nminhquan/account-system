package utils

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

// NewSHAHash generates a new SHA1 hash based on
// a random number of characters.
func NewSHAHash(strs ...string) string {
	hash := sha1.New()
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}

	hash.Write([]byte(sb.String()))
	bs := hash.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
