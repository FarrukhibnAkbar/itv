package utils

import "github.com/google/uuid"

// IsValidUUID validates uuid
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
