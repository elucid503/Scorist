package commands

import (
	"fmt"
	"log"
	"paul/scorist/fetcher"
	"paul/scorist/models"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

const (
	embedColorScheduled = 0x5865F2
	embedColorLive = 0x57F287
	embedColorFinal = 0x95A5A6
	embedColorDefault = 0x2B2D31
)

func ViewScheduleCommand(event *events.ApplicationCommandInteractionCreate) {

	schedule, err := fetcher.FetchToday()

	if err != nil {

		log.Printf("Error fetching schedule: %v", err)

		respondError(event, "Could not load today's schedule.")

		return

	}

	games := fetcher.TodayGames(schedule)

	if len(games) == 0 {

		respondEmbeds(event, discord.NewEmbed().
			WithTitle("MLB Schedule").
			WithDescription("No games are scheduled for today.").
			WithColor(embedColorDefault))

		return

	}

	respondEmbeds(event, buildScheduleEmbeds(scheduleDateLabel(schedule), games)...)

}

func ViewGameCommand(event *events.ApplicationCommandInteractionCreate) {

	gamePk, err := strconv.Atoi(event.ApplicationCommandInteraction.SlashCommandInteractionData().String("game"))

	if err != nil {

		respondError(event, "Invalid game selection.")

		return

	}

	schedule, err := fetcher.FetchToday()

	if err != nil {

		log.Printf("Error fetching schedule: %v", err)

		respondError(event, "Could not load game data.")

		return

	}

	game, ok := fetcher.FindGame(schedule, gamePk)

	if !ok {

		respondError(event, "That game is not on today's schedule.")

		return

	}

	respondEmbeds(event, buildGameEmbed(game))

}

func ViewGameAutocomplete(event *events.AutocompleteInteractionCreate) {

	data := event.AutocompleteInteraction.Data

	if data.CommandName != "view-game" {

		return

	}

	focused := data.Focused()

	if focused.Name != "game" {

		return

	}

	query := ""

	if focused.Type == discord.ApplicationCommandOptionTypeString {

		query = focused.String()

	}

	schedule, err := fetcher.FetchToday()

	if err != nil {

		log.Printf("Error fetching schedule for autocomplete: %v", err)

		return

	}

	choices := gameAutocompleteChoices(fetcher.TodayGames(schedule), query)

	err = event.AutocompleteResult(choices)

	if err != nil {

		log.Printf("Error sending autocomplete: %v", err)

	}

}

func buildScheduleEmbeds(dateLabel string, games []models.Game) []discord.Embed {

	embeds := make([]discord.Embed, 0, (len(games)+24)/25)

	for len(games) > 0 {

		chunk := games

		if len(chunk) > 25 {

			chunk = games[:25]

			games = games[25:]

		} else {

			games = nil

		}

		embed := discord.NewEmbed().
			WithTitle("MLB Schedule").
			WithDescription(dateLabel).
			WithColor(embedColorDefault)

		for _, game := range chunk {

			embed = embed.AddField(gameMatchup(game), gameScheduleValue(game), true)

		}

		embeds = append(embeds, embed)

	}

	return embeds

}

func buildGameEmbed(game models.Game) discord.Embed {

	away := game.Teams.Away
	home := game.Teams.Home

	embed := discord.NewEmbed().
		WithTitle(gameMatchup(game)).
		WithColor(statusEmbedColor(game.Status.DetailedState))

	embed = embed.AddField("Status", game.Status.DetailedState, true)

	if score := gameScoreText(game); score != "" {

		embed = embed.AddField("Score", score, true)

	}

	embed = embed.AddField("Start", formatGameStart(game, "f"), true)

	embed = embed.AddField(away.Team.Name, fmt.Sprintf("%d-%d", away.LeagueRecord.Wins, away.LeagueRecord.Losses), false)
	embed = embed.AddField(home.Team.Name, fmt.Sprintf("%d-%d", home.LeagueRecord.Wins, home.LeagueRecord.Losses), false)

	if game.Venue.Name != "" {

		embed = embed.AddField("Venue", game.Venue.Name, true)

	}

	if game.SeriesDescription != "" {

		embed = embed.AddField("Series", game.SeriesDescription, true)

	}

	return embed

}

func scheduleDateLabel(schedule *models.Schedule) string {

	if schedule != nil && len(schedule.Dates) > 0 {

		if parsed, err := time.ParseInLocation("2006-01-02", schedule.Dates[0].Date, time.UTC); err == nil {

			return discordTimestamp(parsed, "D")

		}

	}

	return discordTimestamp(time.Now().UTC(), "D")

}

func gameScheduleValue(game models.Game) string {

	lines := []string{game.Status.DetailedState}

	if score := gameScoreText(game); score != "" {

		lines = append(lines, score)

	}

	lines = append(lines, formatGameStart(game, "t"))

	return strings.Join(lines, "\n")

}

func gameMatchup(game models.Game) string {

	name := fmt.Sprintf("%s @ %s", game.Teams.Away.Team.Name, game.Teams.Home.Team.Name)

	if len(name) > 256 {

		return name[:253] + "..."

	}

	return name

}

func gameAutocompleteChoices(games []models.Game, query string) []discord.AutocompleteChoice {

	choices := make([]discord.AutocompleteChoice, 0, 25)

	query = strings.ToLower(strings.TrimSpace(query))

	for _, game := range games {

		if query != "" && !gameMatchesQuery(game, query) {

			continue

		}

		choices = append(choices, discord.AutocompleteChoiceString{

			Name: gameChoiceName(game),
			Value: strconv.Itoa(game.GamePk),

		})

		if len(choices) >= 25 {

			break

		}

	}

	return choices

}

func gameChoiceName(game models.Game) string {

	label := fmt.Sprintf("%s @ %s (%s)", game.Teams.Away.Team.Name, game.Teams.Home.Team.Name, game.Status.DetailedState)

	if len(label) > 100 {

		return label[:97] + "..."

	}

	return label

}

func gameMatchesQuery(game models.Game, query string) bool {

	if strings.Contains(strconv.Itoa(game.GamePk), query) {

		return true

	}

	return strings.Contains(strings.ToLower(gameChoiceName(game)), query)

}

func gameScoreText(game models.Game) string {

	state := game.Status.DetailedState

	if state != "In Progress" && state != "Final" && state != "Game Over" {

		return ""

	}

	return fmt.Sprintf("%d – %d", game.Teams.Away.Score, game.Teams.Home.Score)

}

func statusEmbedColor(status string) int {

	switch status {

		case "In Progress":

			return embedColorLive

		case "Final", "Game Over":

			return embedColorFinal

		case "Scheduled", "Pre-Game", "Warmup", "Delayed Start", "Postponed":

			return embedColorScheduled

		default:

			return embedColorDefault

	}

}

func formatGameStart(game models.Game, style string) string {

	if game.Status.StartTimeTBD {

		return "TBD"

	}

	return discordTimestamp(game.GameDate, style)

}

func discordTimestamp(t time.Time, style string) string {

	return fmt.Sprintf("<t:%d:%s>", t.Unix(), style)

}