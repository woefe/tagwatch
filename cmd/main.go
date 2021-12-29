package main

import (
	"log"
	"os"
	"tagwatch"
)

func main() {
	conf, err := tagwatch.LoadConf("tagwatch.example.yml")
	if err != nil {
		log.Fatalln(err)
		return
	}
	switch os.Args[1] {
	case "run":
		_, err := os.Stdout.Write(*tagwatch.MakeFeed(conf))
		if err != nil {
			log.Fatalln(err)
			return
		}
	case "serve":
		tagwatch.Serve(conf)
	}
}
