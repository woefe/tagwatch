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
	"log"
	"net/http"
	"sync"
	"time"
)

var emptyFeedBytes = []byte(EmptyFeed)
var feedXML = &emptyFeedBytes
var feedMutex sync.Mutex
var conf *Conf

func handleFeed(resp http.ResponseWriter, req *http.Request) {
	feedMutex.Lock()
	defer feedMutex.Unlock()
	if req.Method != http.MethodGet {
		log.Println("Invalid Request")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	headers := resp.Header()
	headers["Content-Type"] = []string{"application/xml"}
	_, err := resp.Write(*feedXML)
	if err != nil {
		log.Println(err)
	}
}

func refreshFeed() {
	feedMutex.Lock()
	defer feedMutex.Unlock()
	feedXML = MakeFeed(conf)
}

func backgroundGenerator() {
	lastRun := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	for {
		if lastRun.Add(4 * time.Hour).Before(time.Now()) {
			lastRun = time.Now()
			log.Println("Refreshing feed...")
			refreshFeed()
		}
		time.Sleep(1 * time.Minute)
	}
}

func Serve(config *Conf) {
	conf = config
	log.Println("Starting background feed generator...")
	go backgroundGenerator()

	log.Println("Starting server on", config.Server.Addr, "...")
	http.HandleFunc("/feed.xml", handleFeed)
	log.Fatalln(http.ListenAndServe(config.Server.Addr, nil))
}
