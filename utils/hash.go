package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func HashPassword(password string) string {
	hashedPassHex := md5.Sum([]byte(password))
	hashedPassString := hex.EncodeToString(hashedPassHex[:])
	for i := 0; i < len(hashedPassString); i = i + 2 {
		if hashedPassString[i] == 0 {
			hashedPassString = fmt.Sprintf("%sc%s", hashedPassString[0:i], hashedPassString[i+1:])
		}
	}
	return hashedPassString
}
