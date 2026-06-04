package discord

import (
	"context"
	"fmt"
	"log"
	"paul/scorist/db"
	"paul/scorist/fetcher"
	"paul/scorist/models"
	"time"

	"github.com/disgoorg/disgo/discord"
)

type Notifier struct {

	Store *db.Store

}

func NewNotifier(store *db.Store) *Notifier {

	return &Notifier{Store: store}

}

func (n *Notifier) Handle(event fetcher.ScoreEvent) {

	ctx := context.Background()

	guilds, err := n.Store.ListConfigured(ctx)

	if err != nil {

		log.Printf("Error loading guild settings: %v", err)

		return

	}

	for _, guild := range guilds {

		if event.Final {

			n.schedule(guild, event)

		} else if guild.Preferences.Updates == models.UpdateModeScoreChanges {

			n.schedule(guild, event)

		}

	}

}

func (n *Notifier) schedule(guild *models.Guild, event fetcher.ScoreEvent) {

	delay := guild.Preferences.Delay.Duration()

	g := *guild
	e := event

	go func() {

		time.Sleep(delay)

		if Client == nil {

			return

		}

		_, err := Client.Rest.CreateMessage(g.Channels.Updates, discord.NewMessageCreate().WithContent(formatScoreMessage(e)))

		if err != nil {

			log.Printf("Error posting score update to guild %d: %v", g.ID, err)

		}

	}()

}

func formatScoreMessage(event fetcher.ScoreEvent) string {

	scoreline := fmt.Sprintf("**%s** %d - %d **%s**", event.AwayName, event.AwayScore, event.HomeScore, event.HomeName)

	if event.Final {

		return fmt.Sprintf("🏁 Final: %s", scoreline)

	}

	if event.Inning != "" {

		return fmt.Sprintf("⚾ %s (%s)", scoreline, event.Inning)

	}

	return fmt.Sprintf("⚾ %s", scoreline)

}