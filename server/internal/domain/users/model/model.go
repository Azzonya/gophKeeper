// Package model defines the data structures for managing user accounts,
// including user information, query parameters, and edit operations.
package model

import "time"

// Main represents the core user entity, storing user identification details,
// username, password hash, and timestamps for record creation and updates.
type Main struct {
	UserID       string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GetPars defines parameters for querying specific user records,
// allowing filtering by UserID or Username.
type GetPars struct {
	UserID   string
	Username string
}

// IsValid checks if at least one field in GetPars is populated.
func (m *GetPars) IsValid() bool {
	return m.UserID != "" || m.Username != ""
}

// ListPars defines parameters for listing user records with optional filters,
// supporting filtering by UserIDs, Username, and timestamps.
type ListPars struct {
	UserID        *string
	UserIDs       *[]string
	Username      *string
	CreatedBefore *time.Time
	CreatedAfter  *time.Time
	UpdatedBefore *time.Time
	UpdatedAfter  *time.Time
}

// Edit represents the editable fields for updating an existing user record,
// allowing partial updates to fields like Username, PasswordHash, and timestamps.
type Edit struct {
	UserID       string
	Username     *string
	PasswordHash *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}
