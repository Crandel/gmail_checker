package extractors

import "encoding/xml"

type Feed struct {
	XMLName   xml.Name `xml."feed"`
	Xmlns     string   `xml."xmlns,attr"`
	Version   string   `xml."version,attr"`
	Title     string   `xml."name"`
	Tagline   string   `xml."tagline"`
	Fullcount string   `xml."fullcount"`
	Link      Link     `xml."link"`
}

type Link struct {
	XMLName xml.Name `xml."link"`
	Rel     string   `xml."rel,attr"`
	Href    string   `xml."href,attr"`
	Type    string   `xml."type,attr"`
}
