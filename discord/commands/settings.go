package commands

import (
	"context"
	"fmt"
	"log"
	"paul/scorist/db"
	"paul/scorist/models"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Store *db.Store

func guildFromEvent(event *events.ApplicationCommandInteractionCreate) (*models.Guild, error) {

	guildID := event.GuildID()

	if guildID == nil {

		return nil, errGuildOnly

	}

	return Store.Get(context.Background(), *guildID)

}

func saveGuild(guild *models.Guild) error {

	return Store.Save(context.Background(), guild)

}

func SetChannelCommand(event *events.ApplicationCommandInteractionCreate) {

	guild, err := guildFromEvent(event)

	if err != nil {

		respondError(event, err.Error())

		return

	}

	data := event.ApplicationCommandInteraction.SlashCommandInteractionData()

	option, ok := data.Option("channel")

	if !ok {

		respondError(event, "Channel option is required.")

		return

	}

	channelID := option.Snowflake()

	guild.SetChannel(models.ChannelTypeUpdates, channelID)

	if err := saveGuild(guild); err != nil {

		log.Printf("Error saving guild settings: %v", err)

		respondError(event, "Failed to save settings.")

		return

	}

	channelRef := fmt.Sprintf("<#%s>", channelID)

	if channel, ok := data.OptChannel("channel"); ok {

		channelRef = "#" + channel.Name

	}

	respond(event, fmt.Sprintf("Score updates will be posted in %s.", channelRef))

}

func SetUpdatesCommand(event *events.ApplicationCommandInteractionCreate) {

	guild, err := guildFromEvent(event)

	if err != nil {

		respondError(event, err.Error())

		return

	}

	mode := models.UpdateMode(event.ApplicationCommandInteraction.SlashCommandInteractionData().String("updates"))

	if mode != models.UpdateModeFinalScores && mode != models.UpdateModeScoreChanges {

		respondError(event, "Invalid update mode.")

		return

	}

	guild.SetUpdates(mode)

	if err := saveGuild(guild); err != nil {

		log.Printf("Error saving guild settings: %v", err)

		respondError(event, "Failed to save settings.")

		return

	}

	respond(event, fmt.Sprintf("Update mode set to `%s`.", mode))

}

func SetDelayCommand(event *events.ApplicationCommandInteractionCreate) {

	guild, err := guildFromEvent(event)

	if err != nil {

		respondError(event, err.Error())

		return

	}

	delay := models.DelayMode(event.ApplicationCommandInteraction.SlashCommandInteractionData().String("time"))

	switch delay {

		case models.Delay1Min, models.Delay5Min, models.Delay30Min, models.DelayRandom:

		default:

			respondError(event, "Invalid delay.")

			return

	}

	guild.SetDelay(delay)

	if err := saveGuild(guild); err != nil {

		log.Printf("Error saving guild settings: %v", err)

		respondError(event, "Failed to save settings.")

		return

	}

	respond(event, fmt.Sprintf("Post delay set to `%s`.", delayLabel(delay)))

}

func delayLabel(delay models.DelayMode) string {

	switch delay {

		case models.Delay1Min:

			return "1 min"

		case models.Delay5Min:

			return "5 min"

		case models.Delay30Min:

			return "30 min"

		case models.DelayRandom:

			return "random"

		default:

			return string(delay)

	}

}

func respond(event *events.ApplicationCommandInteractionCreate, content string) {

	err := event.CreateMessage(discord.NewMessageCreate().WithContent(content))

	if err != nil {

		log.Printf("Error responding to command: %v", err)

	}

}

func respondEmbeds(event *events.ApplicationCommandInteractionCreate, embeds ...discord.Embed) {

	err := event.CreateMessage(discord.NewMessageCreate().WithEmbeds(embeds...))

	if err != nil {

		log.Printf("Error responding to command: %v", err)

	}

}

func respondError(event *events.ApplicationCommandInteractionCreate, content string) {

	respond(event, content)

}

var errGuildOnly = guildOnlyError{}

type guildOnlyError struct{}

func (guildOnlyError) Error() string {

	return "These commands can only be used in a server."

}