package commands

import (
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func PingCommand(event *events.ApplicationCommandInteractionCreate) {

	err := event.CreateMessage(discord.NewMessageCreate().WithContent("Pong!"))

	if err != nil {

		log.Printf("Error responding to ping command: %v", err)

	}

}
