package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func HashWithSalt(s, salt string) string {
	concat := fmt.Sprintf("%s_%s", s, salt)
	hashSign := md5.Sum([]byte(concat))
	signStr := hex.EncodeToString(hashSign[:])
	return signStr
}
