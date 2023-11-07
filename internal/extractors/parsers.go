package extractors

import (
	"encoding/xml"
	"log/slog"
)

func ExtractCount(str string) string {
	var feed Feed
	err := xml.Unmarshal([]byte(str), &feed)
	if err != nil {
		slog.Debug("Error during extracting counts", slog.Any("error", err))
	}
	return feed.Fullcount
}
