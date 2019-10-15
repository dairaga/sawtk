package util

import uuid "github.com/google/uuid"

// NewUUID returns a random uuid string (uuid v4).
func NewUUID() string {
	return uuid.New().String()
}

// IsUUID returns id is an UUID string or not.
func IsUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
