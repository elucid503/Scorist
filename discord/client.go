package discord

import (
	"context"
	"paul/scorist/discord/commands"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Client *bot.Client

func Init(token string) error {

	client, err := disgo.New(token)

	if err != nil {

		return err

	}

	return client.OpenGateway(context.Background())

}

func CreateCommands() error {

	_, err := Client.Rest.SetGlobalCommands(Client.ApplicationID, []discord.ApplicationCommandCreate{

		discord.SlashCommandCreate{

			Name: "ping",
			Description: "Use this command to test if the bot is responsive.",

		},

	})

	return err

}

func RegisterEvents() {

	Client.AddEventListeners(bot.NewListenerFunc(func (event *events.ApplicationCommandInteractionCreate) {

		name := event.ApplicationCommandInteraction.SlashCommandInteractionData().CommandName()

		switch name {

			case "ping":

				commands.PingCommand(event)

		}

	}))

}
