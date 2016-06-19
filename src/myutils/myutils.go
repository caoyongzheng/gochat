package myutils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

//UniqueID 生成一个唯一ID值
func UniqueID() (string, error) {
	b := make([]byte, 16)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG.")
	}
	return hex.EncodeToString(b), nil
}
