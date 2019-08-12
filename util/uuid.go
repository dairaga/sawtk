package util

import uuid "github.com/satori/go.uuid"

// NewUUID returns a random uuid string (uuid v4).
func NewUUID() string {
	return uuid.NewV4().String()
}

// IsUUID returns id is an UUID string or not.
func IsUUID(id string) bool {
	_, err := uuid.FromString(id)
	return err == nil
}
