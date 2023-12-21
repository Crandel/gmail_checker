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
	"github.com/Crandel/gmail/internal/mails/gmail"
)

const gmailUrl = "https://mail.google.com"

func main() {
	debug := os.Getenv("DEBUG")
	logLevel := slog.LevelInfo
	showSources := false
	if debug == "1" {
		showSources = true
		logLevel = slog.LevelDebug
	}

	logging.InitLogger(logLevel, showSources)
	createFlag := flag.Bool("create", false, "Create example configuration")
	loadPathFlag := flag.String("load", "", "Load credentials.json")

	flag.Parse()

	if *createFlag {
		config.CreateConfig()
		return
	}
	if loadPathFlag != nil && *loadPathFlag != "" {
		fmt.Printf("The path is: /n%s/n", *loadPathFlag)
		gmail.SaveConfig(*loadPathFlag)
	}
	// Check if domain online
	resp, err := http.Get(gmailUrl)
	ctx := context.Background()

	if err == nil || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently {
		channel := make(chan string)
		defer close(channel)
		listAccounts := config.GetAccounts()
		for _, acc := range listAccounts {
			if acc.MailType == gmail.Type {
				config, err := gmail.GetConfig(acc.ClientID)
				if err != nil {
					slog.Debug("can't get config", slog.Any("error", err))
				}
				// separate all network requests to goroutines

				client := gmail.GetClient(config)

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
					fmt.Println("No labels found.")
					return
				}
				fmt.Println("Labels:")
				for _, l := range r.Labels {
					fmt.Printf("- %s\n", l.Name)
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
}
