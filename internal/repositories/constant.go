package repositories

const (
	//TableNameUsers is a const table name
	TableNameUsers = `users`
	//TableNameBuckets is a const table name
	TableNameBuckets = `buckets`
	//TableNameMovies is a const table name
	TableNameMovies = `movies`
	//TableNameMovieVotes is a const table name
	TableNameMovieVotes = `movie_votes`
	//TableNameMovieGenre is a const table name
	TableNameMovieGenre = `movie_genre`
	//TableNameGenres is a const table name
	TableNameGenres = `genres`
	//TableNameActionHistories is a const table name
	TableNameActionHistories = `action_histories`
)

var (
	//DefaultQueryFindOne is a raw select query
	DefaultQueryFindOne = `SELECT %s FROM %s %s LIMIT 1;`
	//DefaultQueryFinds is a raw static select query
	DefaultQueryFinds = `SELECT %s FROM %s %s;`
)
