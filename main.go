package main

import (
	"fmt"
	"log"
	"os"
)

func usage() {
	log.Println("Invalid arguments!")
	fmt.Println()
	help()
	os.Exit(1)
}

func help() {
	fmt.Println("Usage: tagwatch (run|serve|version|help|-h|--help)")
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	switch os.Args[1] {
	case "-h", "--help", "help":
		help()
		os.Exit(0)
	case "version":
		fmt.Print("tagwatch version ")
		fmt.Println(Version)
		os.Exit(0)
	case "run", "serve":
	default:
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
	}
}
