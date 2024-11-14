package logging

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// InitLogger will init logger globally.
func InitLogger(level slog.Level, show bool) {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.DateTime,
		Level:      level,
		AddSource:  show,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
