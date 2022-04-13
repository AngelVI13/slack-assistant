package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AngelVI13/slack-assistant/slack/handlers"

	"github.com/AngelVI13/slack-assistant/utils"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
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

	// Load Env variables from .dot file
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	debug := false
	if os.Getenv("SL_DEBUG") == "1" {
		debug = true
	}

	// Create a new client to slack
	client := slack.New(token, slack.OptionDebug(debug), slack.OptionAppLevelToken(appToken))

	// Convert simple slack client to socket mode client
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(debug),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	devicesFile := os.Getenv("SL_DEVICES_FILE")
	workersEndpoint := fmt.Sprintf("%s/workers", os.Getenv("SL_TA_ENDPOINT"))
	devicesInfo := utils.GetDevices(devicesFile, workersEndpoint)

	usersFile := os.Getenv("SL_USERS_FILE")
	users := utils.GetUsers(usersFile)
	deviceManager := handlers.DeviceManager{
		DevicesMap:  devicesInfo,
		Users:       users,
		SlackClient: socketClient,
	}

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go deviceManager.ProcessMessageLoop(ctx)

	socketClient.Run()
}
