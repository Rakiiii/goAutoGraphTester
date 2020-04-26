package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var lo logger = logger{state: true}

func main() {
	//reading config
	config, err := readConfig(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	initFilesAndDirs(config)
	condition := NewTestState(config)

	switch {
	case config.TypeOfTest == EDGETEST || config.TypeOfTest == VERTEXTEST:
		if condition, err = startDeffTest(config); err != nil {
			log.Panic(err)
		}
	case config.TypeOfTest == ITTEST:
		if condition, err = startItterationTest(config); err != nil {
			log.Panic(err)
		}
	case config.TypeOfTest == PARSETEST:
		if condition, err = startParseTest(config); err != nil {
			log.Panic(err)
		}
	default:
		log.Panicln("Wrong type of test")
		return
	}

	if config.TypeOfTest != PARSETEST {
		//open file with results
		resFile, err := os.Open(PathToResultFile)
		if err != nil {
			log.Println(err)
			return
		}
		defer resFile.Close()

		//draw graphic
		err = DrawPlotCust(resFile, config, condition.Itterator())
		if err != nil {
			fmt.Println(err)
			return
		}

		if config.MTCFG.DrawDiffGraphic {
			lo.log("draw diff graphic")
			resFile.Seek(0, io.SeekStart)
			err = DrawMarkDiff(resFile, config, condition.Itterator())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if config.MTCFG.DrawDynGraphic {
			resFile.Seek(0, io.SeekStart)
			lo.log("draw progression graphic")
			err = DrawMarkProgression(resFile, config, condition.Itterator())
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if config.ATCFG.DrawDistribGraphic {
			advtimeFile, err := os.Open(PathToAdvTimeFile)
			if err != nil {
				log.Println(err)
				return
			}
			defer advtimeFile.Close()
			lo.log("draw distribution graphic")
			err = DrawAdvtimeDistribution(advtimeFile, config, condition.Itterator())
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
