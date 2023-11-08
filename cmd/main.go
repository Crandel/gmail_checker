package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Crandel/gmail/internal/config"
	"github.com/Crandel/gmail/internal/logging"
	"github.com/Crandel/gmail/internal/mails"
)

func main() {
	debug := os.Getenv("DEBUG")
	logLevel := slog.LevelInfo
	showSources := false
	if debug == "1" {
		showSources = true
		logLevel = slog.LevelDebug
	}

	logging.InitLogger(logLevel, showSources)
	// Check if domain online
	resp, err := http.Get("https://mail.google.com")
	if err == nil || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently {
		channel := make(chan string)
		defer close(channel)
		listAccounts := config.GetAccounts()
		for _, acc := range listAccounts {
			// separate all network requests to goroutines
			go mails.GetMailCount(channel, acc)
		}
		accLen := len(listAccounts)
		counts := make([]string, accLen)
		for i := 0; i < accLen; i++ {
			counts[i] = <-channel
		}
		fmt.Println(strings.Join(counts, ""))
	}
}
