package utils

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/AngelVI13/slack-assistant/config"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func SetupSlackClient(config *config.Config, logWriter io.Writer) *socketmode.Client {
	httpClient := http.DefaultClient
	if config.ProxyUrl != "" {
		// Currently this does not work on Centos7. Investigate if it works on Debian11
		// with default proxy stuff i.e. taking the HTTP_PROXY from environment
		// in that case setting custom proxy should not be needed cause by default the
		// DefaultTransport uses ProxyFromEnvironment
		transport := http.DefaultTransport.(*http.Transport).Clone()

		proxyURL, err := url.Parse(config.ProxyUrl)
		if err != nil {
			log.Fatalf("Couldn't setup proxy.\n%+v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		httpClient = &http.Client{Transport: transport}
	}

	client := slack.New(
		config.SlackAuthToken,
		slack.OptionDebug(config.Debug),
		slack.OptionLog(log.New(logWriter, "client: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(config.SlackAppToken),
		slack.OptionHTTPClient(httpClient),
	)

	// Convert simple slack client to socket mode client
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(config.Debug),
		socketmode.OptionLog(log.New(logWriter, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	return socketClient
}
