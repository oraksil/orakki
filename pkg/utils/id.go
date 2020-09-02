package utils

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
)

func NewId(prefix string) string {
	id, _ := gonanoid.Generate("abcdef", 7)
	if prefix != "" {
		return fmt.Sprintf("%s-%s", prefix, id)
	}
	return id
}
