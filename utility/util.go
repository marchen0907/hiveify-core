package utility

import (
	"crypto/sha512"
	"encoding/hex"
)

// StringSha512 计算字符串的sha512
func StringSha512(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
