package entity

import "time"

type MovieGenre struct {
	ID        string     `json:"id,omitempty" db:"id,omitempty"`
	MovieID   string     `json:"movie_id,omitempty" db:"movie_id,omitempty"`
	GenreID   string     `json:"genre_id,omitempty" db:"genre_id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	IsDeleted bool       `json:"is_deleted,omitempty" db:"is_deleted,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at,omitempty"`
}
