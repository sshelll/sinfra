package util

import (
	"strings"

	"github.com/google/uuid"
)

func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func UUIDUpper() string {
	return strings.ToUpper(UUID())
}

func UUIDWithDash() string {
	return uuid.New().String()
}

func UUIDWithDashUpper() string {
	return strings.ToUpper(UUIDWithDash())
}
