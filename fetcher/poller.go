package fetcher

import (
	"fmt"
	"paul/scorist/models"
	"paul/scorist/utils"
	"slices"
	"time"
)

type Poller struct {

	Interval int

	Schedule models.Schedule
	Games []models.Linescore // in-progress games, updated as they are polled

}


// Constructor

func NewPoller(interval int) *Poller {

	return &Poller{

		Interval: interval,

		Schedule: models.Schedule{},
		Games: make([]models.Linescore, 0),

	}

}

func (p *Poller) Start() {

	p.Poll() // initial poll to populate our data immediately

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

	if len(p.Schedule.Dates) == 0 || p.Schedule.Dates[0].Date != time.Now().Format("2006-01-02") {

		_, err := p.getSchedule()

		if err != nil {

			return err

		}

	}

	return p.Update()

}

func (p *Poller) Update() error {

	// For each game in the schedule, if it's in progress, we want to fetch its linescore and update our games slice

	for _, date := range p.Schedule.Dates {

		for _, game := range date.Games {

			if game.Status.DetailedState == "In Progress" {

				linescore, err := p.getLinescore(game.GamePk)

				if err != nil {

					fmt.Printf("Error fetching linescore for game %d: %v\n", game.GamePk, err)

					continue

				}

				if slices.ContainsFunc(p.Games, func(g models.Linescore) bool { return g.GamePk == game.GamePk }) {

					// update it

					index := slices.IndexFunc(p.Games, func(g models.Linescore) bool { return g.GamePk == game.GamePk })

					p.Games[index] = linescore

					fmt.Printf("Updated game %d in games slice\n", game.GamePk)

				} else {

					// add it

					p.Games = append(p.Games, linescore)

					fmt.Printf("Added game %d to games slice\n", game.GamePk)

				}

			}

		}

	}

	p.Clean()

	return nil

}

func (p *Poller) Clean() {

	// remove games > 24 hrs old

	for i := len(p.Games) - 1; i >= 0; i-- {

		game := p.Games[i]

		if time.Now().Unix() - int64(game.CreatedAt) > 24 * 60 * 60 {

			// remove it

			p.Games = append(p.Games[:i], p.Games[i+1:]...) // removes at index i

			fmt.Printf("Removed game %d from games slice (timed out)\n", game.GamePk)

		}

	}

}

// Inner utilities for the poller

func (p *Poller) getSchedule() (*models.Schedule, error) {

	day := time.Now().Format("2006-01-02") // MLB API expects date in YYYY-MM-DD format

	fmt.Printf("Fetching schedule for %s\n", day)

	var schedule models.Schedule

	err := utils.GetAndDecode(fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s", day), &schedule)

	if err != nil {

		return nil, err

	}

	p.Schedule = schedule

	return &schedule, nil

}

func (p *Poller) getLinescore(gamePk int) (models.Linescore, error) {

	var linescore models.Linescore

	err := utils.GetAndDecode(fmt.Sprintf("https://statsapi.mlb.com/api/v1/game/%d/linescore", gamePk), &linescore)

	linescore.GamePk = gamePk // API doesn't include this, but we want it for internal tracking
	linescore.CreatedAt = int(time.Now().Unix()) // also add a timestamp for when we fetched this

	fmt.Printf("Fetched linescore for game %d: %d-%d\n", gamePk, linescore.Teams.Home.Runs, linescore.Teams.Away.Runs)

	return linescore, err

}
