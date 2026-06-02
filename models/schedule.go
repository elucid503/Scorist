package models

import "time"

// URL: https://statsapi.mlb.com/api/v1/schedule?sportId=1&date={date}

type Schedule struct {

	Copyright string `json:"copyright"`

	TotalItems int `json:"totalItems"`
	TotalEvents int `json:"totalEvents"`
	TotalGames int `json:"totalGames"`

	TotalGamesInProgress int `json:"totalGamesInProgress"`

	Dates []Date `json:"dates"`

}

type Date struct {

	Date string `json:"date"`

	TotalItems int `json:"totalItems"`
	TotalEvents int `json:"totalEvents"`
	TotalGames int `json:"totalGames"`

	TotalGamesInProgress int `json:"totalGamesInProgress"`

	Games []Game `json:"games"`
	Events []interface{} `json:"events"` // not really known

}

type Game struct {

	GamePk int `json:"gamePk"` // primary key
	GameGuid string `json:"gameGuid"`

	Link string `json:"link"`
	GameType string `json:"gameType"`

	Season string `json:"season"`
	GameDate time.Time `json:"gameDate"`
	OfficialDate string `json:"officialDate"`

	Status GameStatus `json:"status"`

	Teams Teams `json:"teams"`
	Venue Venue `json:"venue"`
	Content Content `json:"content"`

	IsTie bool `json:"isTie"`

	// Some misc stats / fields

	GameNumber int `json:"gameNumber"`
	PublicFacing bool `json:"publicFacing"`
	DoubleHeader string `json:"doubleHeader"`
	GamedayType string `json:"gamedayType"`
	Tiebreaker string `json:"tiebreaker"`

	CalendarEventID string `json:"calendarEventID"`
	SeasonDisplay string `json:"seasonDisplay"` // e.g. '2026'
	DayNight string `json:"dayNight"` // 'day' or 'night'

	ScheduledInnings int `json:"scheduledInnings"`
	ReverseHomeAwayStatus bool `json:"reverseHomeAwayStatus"`
	InningBreakLength int `json:"inningBreakLength"`

	GamesInSeries int `json:"gamesInSeries"`
	SeriesGameNumber int `json:"seriesGameNumber"`
	SeriesDescription string `json:"seriesDescription"`

	RecordSource string `json:"recordSource"`

	IfNecessary string `json:"ifNecessary"` // 'Y' or 'N'
	IfNecessaryDescription string `json:"ifNecessaryDescription"`

}

type GameStatus struct {

	AbstractGameState string `json:"abstractGameState"`
	CodedGameState string `json:"codedGameState"`
	DetailedState string `json:"detailedState"` // e.g. 'Scheduled', 'In Progress', 'Final', etc.

	StatusCode string `json:"statusCode"`
	AbstractGameCode  string `json:"abstractGameCode"`

	StartTimeTBD bool `json:"startTimeTBD"`

}

type Teams struct {

	Away TeamRecord `json:"away"`
	Home TeamRecord `json:"home"`

}

type TeamRecord struct {

	Team Team `json:"team"`
	LeagueRecord Record `json:"leagueRecord"`

	Score int `json:"score"`
	IsWinner *bool `json:"isWinner"`

	SplitSquad bool `json:"splitSquad"`

	SeriesNumber int `json:"seriesNumber"`

}

type Team struct {

	ID int `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`

}

type Record struct {

	Wins int `json:"wins"`
	Losses int `json:"losses"`
	Ties int `json:"ties"`

	Pct string `json:"pct"` // winning percentage as a string, e.g. '.644'

}

type Venue struct {

	ID int `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`

}

type Content struct {

	Link string `json:"link"`

}
