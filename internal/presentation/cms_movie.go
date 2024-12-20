package presentation

import "github.com/kholidss/movie-fest-skilltest/internal/entity"

type (
	//ReqCMSCreateMovie holding struct of body request
	ReqCMSCreateMovie struct {
		Title           string   `json:"title,omitempty"`
		GenreIDS        []string `json:"genre_ids,omitempty"`
		Description     string   `json:"description,omitempty"`
		MinutesDuration int      `json:"minutes_duration,omitempty"`
		Artists         []string `json:"artists,omitempty"`
		WatchURL        string   `json:"watch_url,omitempty"`
		FileImage       *File    `json:"image"`
	}

	//RespCMSCreateMovie holding struct of body response
	RespCMSCreateMovie struct {
		ID              string         `json:"id"`
		Title           string         `json:"title,omitempty"`
		Genres          []entity.Genre `json:"genres,omitempty"`
		MinutesDuration int            `json:"minutes_duration,omitempty"`
		Artists         []string       `json:"artists,omitempty"`
		WatchURL        string         `json:"watch_url,omitempty"`
		ImageURL        string         `json:"image_url,omitempty"`
		CreatedBy       CreatedBy      `json:"created_by,omitempty"`
	}
)

type (
	//ReqCMSUpdateMovie holding struct of body request
	ReqCMSUpdateMovie struct {
		MovieID         string   `json:"movie_id,omitempty"`
		Title           string   `json:"title,omitempty"`
		GenreIDS        []string `json:"genre_ids,omitempty"`
		Description     string   `json:"description,omitempty"`
		MinutesDuration int      `json:"minutes_duration,omitempty"`
		Artists         []string `json:"artists,omitempty"`
		WatchURL        string   `json:"watch_url,omitempty"`
		FileImage       *File    `json:"image"`
	}

	//RespCMSUpdateMovie holding struct of body response
	RespCMSUpdateMovie struct {
		Title           string         `json:"title,omitempty"`
		Genres          []entity.Genre `json:"genres,omitempty"`
		MinutesDuration int            `json:"minutes_duration,omitempty"`
		Artists         []string       `json:"artists,omitempty"`
		WatchURL        string         `json:"watch_url,omitempty"`
		ImageURL        string         `json:"image_url,omitempty"`
		CreatedBy       CreatedBy      `json:"created_by,omitempty"`
	}
)

type (
	//ReqCMSMostView holding struct of body request
	ReqCMSMostView struct {
		Value string `json:"value,omitempty"`
		Page  int    `json:"page,omitempty"`
		Limit int    `json:"limit,omitempty"`
	}

	//RespCMSMostView holding struct of body response
	RespCMSMostView struct {
		Genres []entity.Genre          `json:"genres,omitempty"`
		Movies []RespCMSMostViewMoview `json:"movies,omitempty"`
	}

	RespCMSMostViewMoview struct {
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
