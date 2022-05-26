package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/AngelVI13/slack-assistant/config"
	"github.com/AngelVI13/slack-assistant/slack/handlers"
	"github.com/AngelVI13/slack-assistant/utils"
	"github.com/joho/godotenv"
)

func main() {
	// examples taken from: https://towardsdatascience.com/develop-a-slack-bot-using-golang-1025b3e606bc

	// Configure logger
	logFile, err := os.OpenFile("./slack-assistant.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	wrt := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(wrt)

	// Env variables are used to configure slack client, devices & users data
	godotenv.Load(".env")
	config := config.ConfigFromEnv()

	socketClient := utils.SetupSlackClient(config, wrt)
	dataHolder := utils.SetupDataHolder(config)

	slackBot := handlers.SlackBot{
		Data:        dataHolder,
		SlackClient: socketClient,
	}

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go slackBot.ProcessMessageLoop(ctx)

	socketClient.Run()
}
