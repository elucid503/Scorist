package discord

import (
	"context"
	"paul/scorist/discord/commands"
	"paul/scorist/models"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

var Client *bot.Client

func Init(token string) error {

	client, err := disgo.New(token, bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsNonPrivileged)))

	if err != nil {

		return err

	}

	Client = client

	return client.OpenGateway(context.Background())

}

func CreateCommands() error {

	_, err := Client.Rest.SetGlobalCommands(Client.ApplicationID, []discord.ApplicationCommandCreate{

		discord.SlashCommandCreate{

			Name: "ping",
			Description: "Use this command to test if the bot is responsive.",

		},

		discord.SlashCommandCreate{

			Name: "set-channel",
			Description: "Set the channel where score updates are posted.",
			Options: []discord.ApplicationCommandOption{

				discord.ApplicationCommandOptionChannel{

					Name: "channel",
					Description: "The text channel for score updates",
					Required: true,
					ChannelTypes: []discord.ChannelType{discord.ChannelTypeGuildText},

				},

			},

		},

		discord.SlashCommandCreate{

			Name: "set-updates",
			Description: "Choose whether to post live score changes, final scores only, or both.",
			Options: []discord.ApplicationCommandOption{

				discord.ApplicationCommandOptionString{

					Name: "updates",
					Description: "final_scores or score_changes (final scores are always sent)",
					Required: true,
					Choices: []discord.ApplicationCommandOptionChoiceString{

						{Name: "final_scores", Value: string(models.UpdateModeFinalScores)},
						{Name: "score_changes", Value: string(models.UpdateModeScoreChanges)},

					},

				},

			},

		},

		discord.SlashCommandCreate{

			Name: "set-delay",
			Description: "Set how long to wait before posting an update.",
			Options: []discord.ApplicationCommandOption{

				discord.ApplicationCommandOptionString{

					Name: "time",
					Description: "Delay before posting",
					Required: true,
					Choices: []discord.ApplicationCommandOptionChoiceString{

						{Name: "1 min", Value: string(models.Delay1Min)},
						{Name: "5 min", Value: string(models.Delay5Min)},
						{Name: "30 min", Value: string(models.Delay30Min)},
						{Name: "random", Value: string(models.DelayRandom)},

					},

				},

			},

		},

		discord.SlashCommandCreate{

			Name: "view-schedule",
			Description: "View today's MLB schedule.",

		},

		discord.SlashCommandCreate{

			Name: "view-game",
			Description: "View details for a game on today's schedule.",
			Options: []discord.ApplicationCommandOption{

				discord.ApplicationCommandOptionString{

					Name: "game",
					Description: "A game from today's schedule",
					Required: true,
					Autocomplete: true,

				},

			},

		},

	})

	return err

}

func RegisterEvents() {

	Client.AddEventListeners(bot.NewListenerFunc(func(event *events.ApplicationCommandInteractionCreate) {

		name := event.ApplicationCommandInteraction.SlashCommandInteractionData().CommandName()

		switch name {

			case "ping":

				commands.PingCommand(event)

			case "set-channel":

				commands.SetChannelCommand(event)

			case "set-updates":

				commands.SetUpdatesCommand(event)

			case "set-delay":

				commands.SetDelayCommand(event)

			case "view-schedule":

				commands.ViewScheduleCommand(event)

			case "view-game":

				commands.ViewGameCommand(event)

		}

	}))

	Client.AddEventListeners(bot.NewListenerFunc(func(event *events.AutocompleteInteractionCreate) {

		commands.ViewGameAutocomplete(event)

	}))

}