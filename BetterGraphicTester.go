package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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

	graphgeneFlag := true

	switch {
	case config.TypeOfTest == EDGETEST || config.TypeOfTest == VERTEXTEST:
		if condition, err = startDeffTest(config); err != nil {
			log.Panic(err)
		}
	case config.TypeOfTest == ITTEST:
		resString := make([]string, 5)

		resString[0] = strconv.Itoa(config.ITCFG.StartingAmountOfItteration + config.ITCFG.ItterrationDifference*condition.Itterator())

		if config.ITCFG.GraphGeneratorBefavor != "no" {
			if graphgeneFlag {
				//generate graph
				if _, err := startGraphGene(config, condition.Itterator()); err != nil {
					log.Println(err)
					return
				}
				lo.log("graph generated")
			}
		}

		if config.ITCFG.GraphGeneratorBefavor == "once" {
			graphgeneFlag = false
		}

		//copy generated graph and get name
		resString[2], err = copyGraph(config, condition.Itterator())
		if err != nil {
			log.Println(err)
			return
		}

		lo.log("graph coppied")

		//start graph handler
		if err := startGraphHandler(config, condition.Itterator()); err != nil {
			log.Println(err)
			return
		}

		lo.log("graph handler finnished")

		//save result and get name
		resString[3], err = saveGraphHandlerResult(config, condition.Itterator())
		if err != nil {
			log.Println(err)
			return
		}

		lo.log("result coppied")

		lo.log("result add to tab")

		//get result value
		resString[1], err = getResult(config)
		if err != nil {
			log.Println(err)
			return
		}

		lo.log("result value getted")

		//get string with mark
		if config.MTCFG.WithFMMark {
			resString[5], err = getMark(config)
			if err != nil {
				log.Println(err)
				return
			}
			lo.log("mark parsed")
		}

		//get string with advtime if enable
		if config.ATCFG.EnableAdvTime {
			advtime := strconv.Itoa(config.ITCFG.StartingAmountOfItteration + config.ITCFG.ItterrationDifference*condition.Itterator())
			tmp, err := getAdvancedTime(config)
			if err != nil {
				log.Println(err)
				return
			}
			advtime += " " + tmp[:(len(tmp)-1)]
			lo.log("advtime parsed")
			if err = AppendStringToFile(PathToAdvTimeFile, advtime, condition.Itterator()); err != nil {
				log.Println(err)
				return
			}
			lo.log("advtime appended")
		}

		//make res string
		writeRes := ""
		for _, i := range resString {
			writeRes += i
			writeRes += " "
		}

		//write result string
		if err = AppendStringToFile(PathToResultFile, writeRes, condition.Itterator()); err != nil {
			log.Println(err)
			return
		}
	case config.TypeOfTest == PARSETEST:

		resString := make([]string, 2)

		setOfGraphs, err := getSliceOfGrpahs(config)
		if err != nil {
			log.Panicln(err)
			return
		}

		log.Println("set of grpahs parsed")

		for it, graphFile := range *setOfGraphs {
			log.Println("grpah handling starts , grpah name:", graphFile)
			config.PTCFG.File = graphFile
			if err := startGraphHandler(config, it); err != nil {
				log.Panicln(err)
				return
			}

			lo.log("graph handler finnished")

			//save result and get name
			resString[0], err = saveGraphHandlerResult(config, condition.Itterator())
			if err != nil {
				log.Println(err)
				return
			}

			lo.log("result coppied")

			//get result value
			resString[1], err = getResultFromParsed(config)
			if err != nil {
				log.Println(err)
				return
			}
			lo.log("result value getted")

			//get string with mark
			resString[5], err = getMark(config)
			if err != nil {
				log.Println(err)
				return
			}

			if config.MTCFG.WithFMMark {
				lo.log("mark parsed")
			}

			//make res string
			writeRes := ""
			for _, i := range resString {
				writeRes += i
				writeRes += " "
			}

			//write result string
			if err = AppendStringToFile(PathToResultFile, writeRes, condition.Itterator()); err != nil {
				log.Println(err)
				return
			}

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
