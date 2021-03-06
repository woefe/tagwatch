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
		if lastRun.Add(12 * time.Hour).Before(time.Now()) {
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

	log.Println("Starting server...")
	http.HandleFunc("/feed.xml", handleFeed)
	log.Fatalln(http.ListenAndServe(config.Server.Addr, nil))
}
