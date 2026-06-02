package main

import (
	"log"
	"os"
	"paul/scorist/discord"
	"paul/scorist/fetcher"

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

	poller := fetcher.NewPoller(30) // 30s interval

	go poller.Start() // starts as a goroutine

	select {} // block forever

}
