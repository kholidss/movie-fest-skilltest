package entity

type Movie struct {
	ID              string `json:"id,omitempty" db:"id,omitempty"`
	Title           string `json:"title,omitempty" db:"title,omitempty"`
	GenreIDS        string `json:"genre_ids,omitempty" db:"genre_ids,omitempty"`
	Description     string `json:"description,omitempty" db:"description,omitempty"`
	MinutesDuration string `json:"minutes_duration,omitempty" db:"minutes_duration,omitempty"`
	ViewNumber      string `json:"view_number,omitempty" db:"view_number,omitempty"`
	Artist          string `json:"artist,omitempty" db:"artist,omitempty"`
	WatchURL        string `json:"watch_url,omitempty" db:"watch_url,omitempty"`
	CreatedBy       string `json:"created_by,omitempty" db:"created_by,omitempty"`
	DefaultCompleteTimestamp
}
