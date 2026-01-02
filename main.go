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
