package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func HashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])
	return passwordHash
}
