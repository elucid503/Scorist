package main

import (
	"log"
	"os"
	"paul/scorist/discord"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {

		log.Fatal("Error loading .env file")

	}

	token := os.Getenv("TOKEN")

	err = discord.Init(token)

	if err != nil {

		log.Fatal("Error initializing discord client: ", err)

	}

	discord.RegisterEvents()
	discord.CreateCommands()

	select {} // block forever

}
