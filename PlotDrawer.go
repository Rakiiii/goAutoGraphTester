package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	config, err := readConfig(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	file, err := os.Open("ResultTab")
	if err != nil {
		log.Println(err)
		return
	}

	am, _ := strconv.Atoi(os.Args[2])
	if err = DrawPlotCust(file, config, am); err != nil {
		log.Println(err)
	}
}
