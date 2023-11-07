package extractors

import "encoding/xml"

func ExtractCount(str string) string {
	var feed Feed
	xml.Unmarshal([]byte(str), &feed)
	return feed.Fullcount
}
