package consts

var (
	//ActionHistoryCreateMovie is a format action history name
	ActionHistoryCreateMovie = `create movie with id: %s, title: %s`

	//ActionHistoryUpdateMovie is a format action history name
	ActionHistoryUpdateMovie = `update movie with id: %s, title: %s`

	//ActionHistoryVoteMovie is a format action history name
	ActionHistoryVoteMovie = `+1 vote movie with id: %s, title: %s`

	//ActionHistoryUnVoteMovie is a format action history name
	ActionHistoryUnVoteMovie = `-1 vote movie with id: %s, title: %s`
)
