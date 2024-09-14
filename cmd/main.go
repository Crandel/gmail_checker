package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	libGmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/Crandel/gmail/internal/config"
	"github.com/Crandel/gmail/internal/logging"
	"github.com/Crandel/gmail/internal/mails/googlemail"
)

const gmailUrl = "https://mail.google.com"

func main() {
	ctx := context.Background()
	debug := os.Getenv("DEBUG")
	logLevel := slog.LevelInfo
	showSources := false
	if debug == "1" {
		showSources = true
		logLevel = slog.LevelDebug
	}

	logging.InitLogger(logLevel, showSources)
	addUserFlag := flag.Bool("add", false, "Add new user to specific mail client")

	flag.Parse()

	if *addUserFlag {
		config.AddToConfig()
		return
	}

	// Check if domain online
	resp, err := http.Get(gmailUrl)

	if err == nil || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently {
		channel := make(chan string)
		defer close(channel)
		listAccounts := config.GetAccounts()
		for _, acc := range listAccounts {
			if acc.MailType == googlemail.Type {
				config, err := googlemail.GetConfig(acc.ClientID, acc.ClientSecret)
				if err != nil {
					slog.Debug("can't load config from credentials", slog.Any("error", err))
					continue
				}
				// separate all network requests to goroutines
				client, err := googlemail.GetClient(ctx, config)
				if err != nil {
					slog.Debug("can't load client for client id: "+acc.ClientID, slog.Any("error", err))
					continue
				}

				srv, err := libGmail.NewService(ctx, option.WithHTTPClient(client))
				if err != nil {
					slog.Debug("Unable to retrieve Gmail client", slog.Any("error", err))
				}

				user := "me"
				r, err := srv.Users.Labels.List(user).Do()
				if err != nil {
					slog.Debug("Unable to retrieve labels", slog.Any("error", err))
				}
				if len(r.Labels) == 0 {
					slog.Debug("No labels found.")
					return
				}
				fmt.Println("Labels:")
				for _, l := range r.Labels {
					fmt.Printf("%s: %d\n", l.Name, l.MessagesUnread)
				}
			}
		}
		// accLen := len(listAccounts)
		// counts := make([]string, accLen)
		// for i := 0; i < accLen; i++ {
		// 	counts[i] = <-channel
		// }
		// fmt.Println(strings.Join(counts, ""))
	}
	slog.Debug("finish")
}
