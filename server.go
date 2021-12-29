package tagwatch

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var emptyFeed = []byte(`<rss version="2.0">
  <channel>
    <title>Docker registry tags</title>
    <link>https://github.com/woefe/tagwatch</link>
    <description>No tags available</description>
    <generator>` + AgentStr + `</generator>
  </channel>
</rss>`)
var feedXML = &emptyFeed
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
		log.Fatalln(err)
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
