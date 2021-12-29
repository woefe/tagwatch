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
	conf, err := LoadConf("tagwatch.example.yml")
	if err != nil {
		log.Fatalln(err)
	}
	if len(os.Args) != 2 {
		usage()
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
