/*
 * Copyright (c) 2026. Wolfgang Popp
 *
 * This file is part of tagwatch.
 *
 * tagwatch is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * tagwatch is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tagwatch.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/xml"
	"time"
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
	Guid        *Guid
}

type Guid struct {
	XMLName     xml.Name `xml:"guid"`
	Content     string   `xml:",chardata"`
	IsPermaLink bool     `xml:"isPermaLink,attr"`
}

type Channel struct {
	XMLName       xml.Name `xml:"channel"`
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Generator     string   `xml:"generator"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Items         []*Item
}

func NewFeed(title, link, description string) *Feed {
	return &Feed{
		Version: "2.0",
		Channel: &Channel{
			Title:         title,
			Link:          link,
			Description:   description,
			Generator:     AgentStr,
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         make([]*Item, 0, 500),
		},
	}
}

func NewItem(title, link, description, guid string) *Item {
	return &Item{
		Title:       title,
		Link:        link,
		Description: description,
		Guid: &Guid{
			Content:     guid,
			IsPermaLink: false,
		},
	}
}

func (feed *Feed) AppendItems(items ...*Item) {
	feed.Channel.Items = append(feed.Channel.Items, items...)
}

func (feed *Feed) ToXML() ([]byte, error) {
	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return nil, err
	}
	return output, nil
}
