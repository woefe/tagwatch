package tagwatch

import (
	"encoding/xml"
	"log"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Channel
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Guid        string   `xml:"guid"`
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Generator   string   `xml:"generator"`
	Items       []*Item
}

func NewFeed(title, link, description string) *Feed {
	return &Feed{
		Version: "2.0",
		Channel: &Channel{
			Title:       title,
			Link:        link,
			Description: description,
			Generator:   "tagwatch/1.0",
			Items:       make([]*Item, 0, 500),
		},
	}
}

func (feed *Feed) AppendItems(items ...*Item) {
	feed.Channel.Items = append(feed.Channel.Items, items...)
}

func (feed *Feed) ToXML() []byte {
	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	return output
}
