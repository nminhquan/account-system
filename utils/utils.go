package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"os"
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

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func GetHostIPv4() string {
	addrs, err := net.InterfaceAddrs()
	var returnIP string
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				returnIP = ipnet.IP.String()
			}
		}
	}

	if returnIP == "" {
		panic("Cannot get Host IP")
	}

	return returnIP
}
