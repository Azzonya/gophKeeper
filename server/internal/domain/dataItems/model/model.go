// Package model defines the core data structures used in the application
// for storing and managing different types of user data, including credentials,
// text, binary data, and bank card information. The package also provides
// structures for handling query parameters and editing operations.
package model

import "time"

const (
	CredentialsDataType = "login_password"
	TextDataType        = "text"
	BinaryDataType      = "binary"
	BankCardDataType    = "bank_card"
)

// Main represents the core data entity, storing user-specific data,
// including the type, content, metadata, and associated timestamps.
type Main struct {
	ID        string
	UserID    string
	Type      string
	Data      []byte
	Meta      string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetPars defines parameters for querying specific records,
// allowing filtering by ID, UserID, Type, Meta, or URL.
type GetPars struct {
	ID     string
	UserID string
	Type   string
	Meta   string
	URL    string
}

// IsValid checks if at least one field in GetPars is populated.
func (m *GetPars) IsValid() bool {
	return m.ID != "" || m.UserID != "" || m.Type != "" || m.Meta != "" || m.URL != ""
}

// ListPars defines parameters for listing records with optional filters,
// supporting filtering by IDs, UserIDs, type, metadata, URL, and timestamps.
type ListPars struct {
	ID            *string
	IDs           *[]string
	UserID        *string
	UserIDs       *[]string
	Type          *string
	Meta          *string
	URL           *string
	CreatedBefore *time.Time
	CreatedAfter  *time.Time
	UpdatedBefore *time.Time
	UpdatedAfter  *time.Time
}

// Edit represents the editable fields for updating an existing record,
// allowing partial updates to fields like Type, Data, Meta, and timestamps.
type Edit struct {
	ID        string
	UserID    *string
	Type      *string
	Data      *[]byte
	Meta      *string
	URL       *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
