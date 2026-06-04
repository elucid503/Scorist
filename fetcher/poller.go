package fetcher

import (
	"fmt"
	"paul/scorist/models"
	"paul/scorist/utils"
	"slices"
	"time"
)

type ScoreHandler func(ScoreEvent)

type Poller struct {

	Interval int

	Schedule models.Schedule
	Games []models.Linescore

	onScore ScoreHandler

	announcedFinals map[int]bool

}

func NewPoller(interval int, onScore ScoreHandler) *Poller {

	return &Poller{

		Interval: interval,

		Schedule: models.Schedule{},
		Games: make([]models.Linescore, 0),

		onScore: onScore,
		announcedFinals: make(map[int]bool),

	}

}

func (p *Poller) Start() {

	p.Poll()

	ticker := time.NewTicker(time.Duration(p.Interval) * time.Second)

	for range ticker.C {

		err := p.Poll()

		if err != nil {

			fmt.Printf("Error polling schedule: %v\n", err)

		}

	}

}

func (p *Poller) Poll() error {

	today := time.Now().Format("2006-01-02")

	if len(p.Schedule.Dates) == 0 || p.Schedule.Dates[0].Date != today {

		p.Games = make([]models.Linescore, 0)
		p.announcedFinals = make(map[int]bool)

		_, err := p.getSchedule()

		if err != nil {

			return err

		}

	}

	return p.Update()

}

func (p *Poller) Update() error {

	for _, date := range p.Schedule.Dates {

		for _, game := range date.Games {

			awayName := game.Teams.Away.Team.Name
			homeName := game.Teams.Home.Team.Name

			switch game.Status.DetailedState {

				case "In Progress":

					linescore, err := p.getLinescore(game.GamePk)

					if err != nil {

						fmt.Printf("Error fetching linescore for game %d: %v\n", game.GamePk, err)

						continue

					}

					prev, tracked := p.findGame(game.GamePk)

					if tracked && scoreChanged(prev, linescore) {

						p.emit(ScoreEvent{

							GamePk: game.GamePk,
							AwayName: awayName,
							HomeName: homeName,
							AwayScore: linescore.Teams.Away.Runs,
							HomeScore: linescore.Teams.Home.Runs,
							Inning: formatInning(linescore),
						})

					}

					p.upsertGame(linescore)

				case "Final":

					if p.announcedFinals[game.GamePk] {

						continue

					}

					awayScore := game.Teams.Away.Score
					homeScore := game.Teams.Home.Score

					if linescore, err := p.getLinescore(game.GamePk); err == nil {

						awayScore = linescore.Teams.Away.Runs
						homeScore = linescore.Teams.Home.Runs

					}

					p.emit(ScoreEvent{

						GamePk: game.GamePk,
						AwayName: awayName,
						HomeName: homeName,
						AwayScore: awayScore,
						HomeScore: homeScore,
						Final: true,
					})

					p.announcedFinals[game.GamePk] = true
					p.removeGame(game.GamePk)

			}

		}

	}

	p.Clean()

	return nil

}

func (p *Poller) emit(event ScoreEvent) {

	if p.onScore == nil {

		return

	}

	p.onScore(event)

}

func scoreChanged(prev, next models.Linescore) bool {

	return prev.Teams.Home.Runs != next.Teams.Home.Runs || prev.Teams.Away.Runs != next.Teams.Away.Runs

}

func formatInning(linescore models.Linescore) string {

	if linescore.InningHalf != "" && linescore.CurrentInningOrdinal != "" {

		half := "Top"

		if !linescore.IsTopInning {

			half = "Bot"

		}

		return fmt.Sprintf("%s %s", half, linescore.CurrentInningOrdinal)

	}

	return linescore.InningState

}

func (p *Poller) findGame(gamePk int) (models.Linescore, bool) {

	index := slices.IndexFunc(p.Games, func(g models.Linescore) bool { return g.GamePk == gamePk })

	if index < 0 {

		return models.Linescore{}, false

	}

	return p.Games[index], true

}

func (p *Poller) upsertGame(linescore models.Linescore) {

	if slices.ContainsFunc(p.Games, func(g models.Linescore) bool { return g.GamePk == linescore.GamePk }) {

		index := slices.IndexFunc(p.Games, func(g models.Linescore) bool { return g.GamePk == linescore.GamePk })

		p.Games[index] = linescore

		fmt.Printf("Updated game %d in games slice\n", linescore.GamePk)

	} else {

		p.Games = append(p.Games, linescore)

		fmt.Printf("Added game %d to games slice\n", linescore.GamePk)

	}

}

func (p *Poller) removeGame(gamePk int) {

	index := slices.IndexFunc(p.Games, func(g models.Linescore) bool { return g.GamePk == gamePk })

	if index < 0 {

		return

	}

	p.Games = append(p.Games[:index], p.Games[index+1:]...)

	fmt.Printf("Removed game %d from games slice\n", gamePk)

}

func (p *Poller) Clean() {

	for i := len(p.Games) - 1; i >= 0; i-- {

		game := p.Games[i]

		if time.Now().Unix()-int64(game.CreatedAt) > 24*60*60 {

			p.Games = append(p.Games[:i], p.Games[i+1:]...)

			fmt.Printf("Removed game %d from games slice (timed out)\n", game.GamePk)

		}

	}

}

func (p *Poller) getSchedule() (*models.Schedule, error) {

	fmt.Printf("Fetching schedule for %s\n", time.Now().Format("2006-01-02"))

	schedule, err := FetchToday()

	if err != nil {

		return nil, err

	}

	p.Schedule = *schedule

	return schedule, nil

}

func (p *Poller) getLinescore(gamePk int) (models.Linescore, error) {

	var linescore models.Linescore

	err := utils.GetAndDecode(fmt.Sprintf("https://statsapi.mlb.com/api/v1/game/%d/linescore", gamePk), &linescore)

	linescore.GamePk = gamePk
	linescore.CreatedAt = int(time.Now().Unix())

	fmt.Printf("Fetched linescore for game %d: %d-%d\n", gamePk, linescore.Teams.Away.Runs, linescore.Teams.Home.Runs)

	return linescore, err

}