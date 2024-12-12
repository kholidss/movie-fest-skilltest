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

type (
	//ReqPublicMovieSearch holding struct of body request
	ReqPublicMovieSearch struct {
		LikeTitle            string `query:"like[title],omitempty" json:"like[title],omitempty"`
		LikeDescription      string `query:"like[description],omitempty" json:"like[description],omitempty"`
		LikeArtist           string `query:"like[artists],omitempty" json:"like[artists],omitempty"`
		LikeWatchURL         string `query:"like[watch_url],omitempty" json:"like[watch_url],omitempty"`
		EqualGenreID         string `query:"equal[genre_id],omitempty" json:"equal[genre_id],omitempty"`
		EqualMinutesDuration string `query:"equal[minutes_duration],omitempty" json:"equal[minutes_duration],omitempty"`

		Page  int `query:"page,omitempty" json:"page,omitempty"`
		Limit int `query:"limit,omitempty" json:"limit,omitempty"`
	}

	//RespPublicMovieSearch holding struct of body response
	RespPublicMovieSearch struct {
		ID              string         `json:"id,omitempty"`
		Title           string         `json:"title,omitempty"`
		Genres          []entity.Genre `json:"genres,omitempty"`
		MinutesDuration int            `json:"minutes_duration,omitempty"`
		ViewNumber      int            `json:"view_number,omitempty"`
		Artists         []string       `json:"artists,omitempty"`
		WatchURL        string         `json:"watch_url,omitempty"`
		ImageURL        string         `json:"image_url,omitempty"`
	}
)
