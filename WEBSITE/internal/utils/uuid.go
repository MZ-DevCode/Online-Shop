package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	fullUUID := uuid.New().String()
	uuid := strings.ReplaceAll(fullUUID, "-", "")
	return uuid[:11]
}
