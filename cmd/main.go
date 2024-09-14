package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Crandel/gmail/internal/accounts"
	"github.com/Crandel/gmail/internal/config"
	"github.com/Crandel/gmail/internal/logging"
	"github.com/Crandel/gmail/internal/mails/googlemail"
)

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

	channel := make(chan string)
	defer close(channel)
	listAccounts := config.GetAccounts()
	for _, acc := range listAccounts {
		if acc.MailType == accounts.Gmail {
			if googlemail.CheckOnline() {
				go googlemail.GetGMailCount(ctx, channel, acc)
			}
		}

	}

	accLen := len(listAccounts)
	counts := make([]string, accLen)
	for i := 0; i < accLen; i++ {
		counts[i] = <-channel
	}
	fmt.Println(strings.Join(counts, ""))
}
