package presentation

type (
	ReqUserMovieVote struct {
		MovieID string `json:"movie_id,omitempty"`
		Value   string `json:"value,omitempty"`
	}

	RespUserMovieVote struct {
		ID    string `json:"id,omitempty"`
		Title string `json:"title,omitempty"`
	}
)
