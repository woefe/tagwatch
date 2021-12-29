package main

import (
	"fmt"
	"log"
	"os"
	"tagwatch"
)

func usage() {
	log.Fatalln("Invalid arguments!\n\nUsage: tagwatch (run|serve|version)")
}

func main() {
	conf, err := tagwatch.LoadConf("tagwatch.example.yml")
	if err != nil {
		log.Fatalln(err)
	}
	if len(os.Args) != 2 {
		usage()
	}
	switch os.Args[1] {
	case "run":
		_, err := os.Stdout.Write(*tagwatch.MakeFeed(conf))
		if err != nil {
			log.Fatalln(err)
		}
	case "serve":
		tagwatch.Serve(conf)
	case "version":
		fmt.Print("tagwatch version ")
		fmt.Println(tagwatch.Version)
	default:
		usage()
	}
}
