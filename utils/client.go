package utils

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func SetupSlackClient(logWriter io.Writer) *socketmode.Client {
	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	debug := false
	if os.Getenv("SL_DEBUG") == "1" {
		debug = true
	}

	proxy := os.Getenv("SL_PROXY")
	httpClient := http.DefaultClient
	if proxy != "" {
		transport := http.DefaultTransport.(*http.Transport).Clone()

		proxyURL, err := url.Parse(proxy)
		if err != nil {
			log.Fatalf("Couldn't setup proxy.\n%+v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		httpClient = &http.Client{Transport: transport}
	}

	client := slack.New(
		token,
		slack.OptionDebug(debug),
		slack.OptionLog(log.New(logWriter, "client: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
		slack.OptionHTTPClient(httpClient),
	)

	// Convert simple slack client to socket mode client
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(debug),
		socketmode.OptionLog(log.New(logWriter, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	return socketClient
}
