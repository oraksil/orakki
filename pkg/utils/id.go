package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

func NewId(prefix string) string {
	id, _ := gonanoid.Generate("abcdef", 7)
	if prefix != "" {
		return fmt.Sprintf("%s-%s", prefix, id)
	}
	return id
}

func NewTurnAuth(userId, turnSecretKey string, turnTTL int) (string, string) {
	timestamp := time.Now().Unix() + int64(turnTTL)
	username := fmt.Sprintf("%d:%s", timestamp, userId)

	h := hmac.New(sha1.New, []byte(turnSecretKey))
	h.Write([]byte(username))

	password := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return username, password
}
