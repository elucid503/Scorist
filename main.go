package main

import (
	"context"
	"log"
	"os"
	"paul/scorist/db"
	"paul/scorist/discord"
	"paul/scorist/discord/commands"
	"paul/scorist/fetcher"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {

		log.Fatal("Error loading .env file")

	}

	mongoURL := os.Getenv("MONGO_URL")

	if mongoURL == "" {

		log.Fatal("MONGO_URL is not set")

	}

	store, err := db.Connect(context.Background(), mongoURL)

	if err != nil {

		log.Fatal("Error connecting to MongoDB: ", err)

	}

	defer store.Close(context.Background())

	commands.Store = store

	token := os.Getenv("TOKEN")

	err = discord.Init(token)

	if err != nil {

		log.Fatal("Error initializing discord client: ", err)

	}

	discord.RegisterEvents()

	err = discord.CreateCommands()

	if err != nil {

		log.Fatal("Error registering commands: ", err)

	}

	notifier := discord.NewNotifier(store)
	poller := fetcher.NewPoller(30, notifier.Handle)

	go poller.Start()

	select {}

}