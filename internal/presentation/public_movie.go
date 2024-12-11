package presentation

import "github.com/kholidss/movie-fest-skilltest/internal/entity"

type (
	ReqPublicMovieList struct {
		Page  int `json:"page,omitempty"`
		Limit int `json:"limit,omitempty"`
	}

	RespPublicMovieList struct {
		ID              string         `json:"id,omitempty"`
		Title           string         `json:"title,omitempty"`
		Genres          []entity.Genre `json:"genres,omitempty"`
		MinutesDuration int            `json:"minutes_duration,omitempty"`
		ViewNumber      int            `json:"view_number,omitempty"`
		Artist          string         `json:"artist,omitempty"`
		WatchURL        string         `json:"watch_url,omitempty"`
		ImageURL        string         `json:"image_url,omitempty"`
	}
)

type (
	ReqPublicTrackMovieViewer struct {
		MovieID string `json:"movie_id,omitempty"`
	}

	RespPublicTrackMovieViewer struct {
		ID    string `json:"id,omitempty"`
		Title string `json:"title,omitempty"`
	}
)
