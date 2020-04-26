package main

import (
	"log"
	"os"
)

//creates all files and direcoryes that are needed
func initFilesAndDirs(config *TestConfig) {
	//create dir for result of graph handler
	if err := createResDir(config); err != nil {
		log.Panic(err)
		return
	}

	lo.log("res dir created")

	//create dir for generated graphs
	if err := createGraphDir(config); err != nil {
		log.Panic(err)
		return
	}

	lo.log("graph dir created")

	//create file for results
	resFile, err := os.Create(PathToResultFile)
	if err != nil {
		resFile.Close()
		log.Panic(err)
	}
	resFile.Close()

	lo.log("result file created")

	//create file for results
	advtimeFile, err := os.Create(PathToAdvTimeFile)
	if err != nil {
		advtimeFile.Close()
		log.Panic(err)
	}
	advtimeFile.Close()
}
