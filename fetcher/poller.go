package fetcher

import (
	"fmt"
	"paul/scorist/models"
	"paul/scorist/utils"
	"time"
)

type Poller struct {

	Interval int

	CurrentSchedule models.Schedule
	Games *Games

}

type Games struct {

	Ongoing map[int]models.Linescore
	Final map[int]models.Linescore

}

// Constructor

func NewPoller(interval int) *Poller {

	return &Poller{

		Interval: interval,

		CurrentSchedule: models.Schedule{},

		Games: &Games{

			Ongoing: make(map[int]models.Linescore),
			Final: make(map[int]models.Linescore),

		},

	}

}

func (p *Poller) Start() {

	ticker := time.NewTicker(time.Duration(p.Interval) * time.Second)

	for range ticker.C {

		// Fires every p.Interval seconds, we will poll the schedule and update our games

		err := p.Poll()

		if err != nil {

			fmt.Printf("Error polling schedule: %v\n", err)

		}

	}

}

func (p *Poller) Poll() error {

	// First we must get the schedule for the current day if needed

	if p.CurrentSchedule.Dates[0].Date != time.Now().Format("2006-01-02") {

		_, err := p.getSchedule()

		if err != nil {

			return err

		}

	}

	return p.Update()

}

func (p *Poller) Update() error {

	for _, date := range p.CurrentSchedule.Dates {

		for _, game := range date.Games {

			currentState := game.Status.DetailedState

			switch currentState {

				case "In Progress":

					// we must always get the linescore for in progress games, as they are always changing

					linescore, err := p.getLinescore(game.GamePk)

					if err != nil {

						return err

					}

					p.Games.Ongoing[game.GamePk] = linescore // replace the linescore for this game

				case "Final":

					// we only need to get the linescore for final games if we haven't already gotten it, as it won't change

					if _, ok := p.Games.Final[game.GamePk]; !ok {

						linescore, err := p.getLinescore(game.GamePk)

						if err != nil {

							return err

						}

						p.Games.Final[game.GamePk] = linescore

					}

			}

		}

	}

	return nil

}

// Inner utilities for the poller

func (p *Poller) getSchedule() (*models.Schedule, error) {

	day := time.Now().Format("2006-01-02") // MLB API expects date in YYYY-MM-DD format

	var schedule models.Schedule

	err := utils.GetAndDecode(fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s", day), &schedule)

	if err != nil {

		return nil, err

	}

	p.CurrentSchedule = schedule

	return &schedule, nil

}

func (p *Poller) getLinescore(gamePk int) (models.Linescore, error) {

	var linescore models.Linescore

	err := utils.GetAndDecode(fmt.Sprintf("https://statsapi.mlb.com/api/v1/game/%d/linescore", gamePk), &linescore)

	return linescore, err

}
