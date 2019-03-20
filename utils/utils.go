package utils

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
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

func GetCurrentTimeInMillis() int64 {
	unixNano := time.Now().UnixNano()
	umillisec := unixNano / 1000000
	return umillisec
}

func GenXid() string {
	return xid.New().String()
}
