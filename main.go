package main

import (
	"fmt"
	"log"
	"os"
)

func usage() {
	log.Fatalln("Invalid arguments!\n\nUsage: tagwatch (run|serve|version)")
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	confPath, exists := os.LookupEnv("TAGWATCH_CONF")
	if !exists {
		confPath = "tagwatch.yml"
	}
	conf, err := LoadConf(confPath)
	if err != nil {
		log.Fatalln(err)
	}

	switch os.Args[1] {
	case "run":
		_, err := os.Stdout.Write(*MakeFeed(conf))
		if err != nil {
			log.Fatalln(err)
		}
	case "serve":
		Serve(conf)
	case "version":
		fmt.Print("tagwatch version ")
		fmt.Println(Version)
	default:
		usage()
	}
}
