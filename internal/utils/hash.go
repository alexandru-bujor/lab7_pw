package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashPasswordSHA256(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func CheckPasswordSHA256(password, stored string) bool {
	fmt.Println(HashPasswordSHA256(password))
	return HashPasswordSHA256(password) == stored
}
