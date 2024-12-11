package entity

import (
	"encoding/json"
	"time"
)

type Bucket struct {
	ID             string          `json:"id,omitempty" db:"id,omitempty"`
	Filename       string          `json:"file_name,omitempty" db:"file_name,omitempty"`
	IdentifierID   string          `json:"identifier_id,omitempty" db:"identifier_id,omitempty"`
	IdentifierName string          `json:"identifier_name,omitempty" db:"identifier_name,omitempty"`
	Mimetype       string          `json:"mimetype,omitempty" db:"mimetype,omitempty"`
	Provider       string          `json:"provider,omitempty" db:"provider,omitempty"`
	URL            string          `json:"url,omitempty" db:"url,omitempty"`
	Path           string          `json:"path,omitempty" db:"path,omitempty"`
	SupportData    json.RawMessage `json:"support_data,omitempty" db:"support_data,omitempty"`
	CreatedAt      time.Time       `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	IsDeleted      bool            `json:"is_deleted,omitempty" db:"is_deleted,omitempty"`
	DeletedAt      *time.Time      `json:"deleted_at,omitempty" db:"deleted_at,omitempty"`
}
