package internal

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func GenerateHash(s ...string) string {
	h := sha256.New()
	h.Write(
		[]byte(
			strings.Join(s, ":"),
		),
	)
	return fmt.Sprintf("%x", h.Sum(nil))
}
