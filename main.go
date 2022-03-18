package main

import (
	"context"
	"github.com/AngelVI13/slack-assistant/handlers"
	"github.com/AngelVI13/slack-assistant/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// examples taken from: https://towardsdatascience.com/develop-a-slack-bot-using-golang-1025b3e606bc

	// Load Env variables from .dot file
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	// Create a new client to slack
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	// Convert simple slack client to socket mode client
	socketClient := socketmode.New(
		client,
		// TODO: Improve logging and load debug flag from cli arg or .env
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	devicesFile := os.Getenv("SL_DEVICES_FILE")
	workersEndpoint := os.Getenv("SL_WORKERS_ENDPOINT")
	devicesInfo := utils.GetDevices(devicesFile, workersEndpoint)

	usersFile := os.Getenv("SL_USERS_FILE")
	users := utils.GetUsers(usersFile)
	deviceManager := handlers.DeviceManager{devicesInfo, users, socketClient}

	// TODO: check if multiple users can interact with the bot at the same time !
	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go deviceManager.ProcessMessageLoop(ctx)

	socketClient.Run()
}
