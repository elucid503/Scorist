package fetcher

type ScoreEvent struct {

	GamePk int

	AwayName string
	HomeName string

	AwayScore int
	HomeScore int

	Inning string // empty for final posts

	Final bool

}