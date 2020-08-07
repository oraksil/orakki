package utils

import (
	"fmt"
	"os"
	"strconv"

	gonanoid "github.com/matoous/go-nanoid"
)

func NewId(prefix string) string {
	id, _ := gonanoid.Generate("abcdef", 7)
	if prefix != "" {
		return fmt.Sprintf("%s-%s", prefix, id)
	}
	return id
}

func GetStrEnv(key string, defValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defValue
	}
	return val
}

func GetIntEnv(key string, defValue int) int {
	val := GetStrEnv(key, "")
	ret, err := strconv.Atoi(val)
	if err != nil {
		return defValue
	}
	return ret
}

func GetBoolEnv(key string, defValue bool) bool {
	val := GetStrEnv(key, "")
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return defValue
	}
	return ret
}
