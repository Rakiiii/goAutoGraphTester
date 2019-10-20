package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var PathToResultFile string = "ResultTab"

func main() {
	//reading config
	config, err := readConfig(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	var lo logger = logger{config.OutputFlag}

	//create dir for result of graph handler
	if err := createResDir(config); err != nil {
		log.Println(err)
		return
	}

	lo.log("res dir created")

	//create dir for generated graphs
	if err := createGraphDir(config); err != nil {
		log.Println(err)
		return
	}

	lo.log("graph dir created")

	//create file for results
	resFile, err := os.Create(PathToResultFile)
	if err != nil {
		fmt.Println(err)
		resFile.Close()
		return
	}
	resFile.Close()

	lo.log("result file created")

	var prevTime time.Duration
	var itterator int = 0
	timeFlag := true
	itterationFlag := true

	for timeFlag && itterationFlag {

		resString := make([]string, 5)

		amV, amE := countAmount(config, itterator)
		resString[0] = strconv.Itoa(amV)
		resString[1] = strconv.Itoa(amE)

		//generate graph
		if err := startGraphGene(config, itterator); err != nil {
			log.Println(err)
			return
		}

		lo.log("graph generated")

		//copy generated graph and get name
		resString[2], err = copyGraph(config, itterator)
		if err != nil {
			log.Println(err)
			return
		}

		lo.log("graph coppied")

		//start graph handler
		if err := startGraphHandler(config, itterator); err != nil {
			log.Println(err)
			return
		}

		lo.log("graph handler finnished")

		//save result and get name
		resString[3], err = saveGraphHandlerResult(config, itterator)
		if err != nil {
			log.Println(err)
			return
		}

		lo.log("result coppied")

		//get string with time
		resString[4], err = getTime(config)
		if err != nil {
			log.Println(err)
			return
		}

		//get time
		prevTime, err = time.ParseDuration(resString[4])
		if err != nil {
			log.Println(err)
			return
		}

		lo.log("time parsed")

		//make res string
		writeRes := ""
		for _, i := range resString {
			writeRes += i
			writeRes += " "
		}

		//write result string
		if err = AppendStringToFile(PathToResultFile, writeRes, itterator); err != nil {
			log.Println(err)
			return
		}

		lo.log("result added to file")

		//stop condition
		if (config.TypeOfStopCondition == "timestop" || config.TypeOfStopCondition == "mixed") && prevTime.Milliseconds() > config.MaxTimeForItteration {
			timeFlag = false
		}

		if (config.TypeOfStopCondition == "itstop" || config.TypeOfStopCondition == "mixed") && config.AmountOfItterations <= itterator {
			itterationFlag = false
		}

		itterator++
	}

	//open file with results
	resFile, err = os.Open(PathToResultFile)
	if err != nil {
		log.Println(err)
		return
	}

	//draw graphic
	err = DrawPlotCust(resFile, config, itterator)
	if err != nil {
		fmt.Println(err)
		return
	}
}
