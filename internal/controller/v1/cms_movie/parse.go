package cmsmovie

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/movie-fest-skilltest/internal/presentation"
	"github.com/kholidss/movie-fest-skilltest/pkg/util"
	"strconv"
)

func (cx *cmsMovieCreate) parse(c *fiber.Ctx) (presentation.ReqCMSCreateMovie, error) {
	var (
		req presentation.ReqCMSCreateMovie
		err error
	)

	req.Title = c.FormValue("title")
	req.Description = c.FormValue("description")
	req.WatchURL = c.FormValue("watch_url")

	req.MinutesDuration, err = strconv.Atoi(func() string {
		if c.FormValue("minutes_duration") == "" {
			return "0"
		}
		return c.FormValue("minutes_duration")
	}())
	if err != nil {
		return req, errors.New("minutes_duration must be a number")
	}

	// Parse the genre_ids JSON string into the slice
	if c.FormValue("genre_ids") != "" {
		err = json.Unmarshal([]byte(c.FormValue("genre_ids")), &req.GenreIDS)
		if err != nil {
			return req, err
		}
	}

	// Parse the artist JSON string into the slice
	if c.FormValue("artists") != "" {
		err = json.Unmarshal([]byte(c.FormValue("artists")), &req.Artists)
		if err != nil {
			return req, err
		}
	}

	// Image File
	req.FileImage, err = util.FiberParseFile(c, "image")
	if err != nil {
		return req, err
	}

	return req, nil
}
