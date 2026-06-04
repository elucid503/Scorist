package fetcher

import (
	"fmt"
	"paul/scorist/models"
	"paul/scorist/utils"
	"time"
)

func FetchToday() (*models.Schedule, error) {

	day := time.Now().Format("2006-01-02")

	var schedule models.Schedule

	err := utils.GetAndDecode(fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s", day), &schedule)

	if err != nil {

		return nil, err

	}

	return &schedule, nil

}

func TodayGames(schedule *models.Schedule) []models.Game {

	if schedule == nil || len(schedule.Dates) == 0 {

		return nil

	}

	return schedule.Dates[0].Games

}

func FindGame(schedule *models.Schedule, gamePk int) (models.Game, bool) {

	for _, game := range TodayGames(schedule) {

		if game.GamePk == gamePk {

			return game, true

		}

	}

	return models.Game{}, false

}