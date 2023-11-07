package extractors

import "encoding/xml"

type Feed struct {
	Feed      xml.Name `xml:"feed"`
	Xmlns     string   `xml:"xmlns,attr"`
	Version   string   `xml:"version,attr"`
	Name      string   `xml:"name"`
	Tagline   string   `xml:"tagline"`
	Fullcount string   `xml:"fullcount"`
	Link      Link     `xml:"link"`
}

type Link struct {
	Link xml.Name `xml:"link"`
	Rel  string   `xml:"rel,attr"`
	Href string   `xml:"href,attr"`
	Type string   `xml:"type,attr"`
}
