package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	libGmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/Crandel/gmail/internal/config"
	"github.com/Crandel/gmail/internal/logging"
	"github.com/Crandel/gmail/internal/mails/googlemail"
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
	createFlag := flag.Bool("config", false, "Update configuration with new mail client")
	loadPathFlag := flag.String("load", "", "Load credentials.json")
	addUserFlag := flag.Bool("add", false, "Add new user to specific mail client")

	flag.Parse()

	if *createFlag {
		config.AddToConfig()
		return
	}

	if loadPathFlag != nil && *loadPathFlag != "" {
		fmt.Printf("The path is: /n %s /n", *loadPathFlag)
		err := googlemail.SaveConfig(*loadPathFlag)
		if err != nil {
			log.Fatal("Can't save config")
		}
	}

	if *addUserFlag {
		addNewUser()
	}

	// Check if domain online
	resp, err := http.Get(gmailUrl)
	ctx := context.Background()

	if err == nil || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently {
		channel := make(chan string)
		defer close(channel)
		listAccounts := config.GetAccounts()
		for _, acc := range listAccounts {
			if acc.MailType == googlemail.Type {
				config, err := googlemail.GetConfig(acc.ClientID)
				if err != nil {
					slog.Debug("can't get config", slog.Any("error", err))
				}
				// separate all network requests to goroutines

				client, err := googlemail.GetClient(config)
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

func addNewUser() {
	// Type necessary information
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please add client id responsible to fetch information for this user")
	clientID, err := reader.ReadString('\n')
	clientID = strings.Trim(clientID, "\n")
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}
	acc := config.GetAccount(clientID)
	if acc == nil {
		fmt.Println("There are no account with this client id. Please try another client id")
		return
	}
	config, err := googlemail.GetConfig(clientID)
	if err != nil {
		fmt.Println("can't get config for this client id. Please try another client id")
		return
	}

	fmt.Println("Please add user alias")
	alias, err := reader.ReadString('\n')
	alias = strings.Trim(alias, "\n")
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}
	err = googlemail.SaveToken(config, alias)
	if err != nil {
		fmt.Println("error saving token to keyring", err)
	}
}
