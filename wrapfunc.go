package main

import (
	"log"
	"os"
	"strconv"
	"time"
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

func startDeffTest(config *TestConfig) (*TestState, error) {
	condition := NewTestState(config)
	ittime, _ := time.ParseDuration("0ms")
	size := *NewGraphSize()
	var err error

	for condition.isContinue(ittime, size.AE, size.AV) {

		lo.log("start vertex of edges test")

		resString := make([]string, 7)

		amV, amE := countAmount(config, condition.Itterator())
		resString[0] = strconv.Itoa(amV)
		resString[1] = strconv.Itoa(amE)

		log.Println("start graph generation")
		//generate graph
		if size, err = startGraphGene(config, condition.Itterator()); err != nil {
			return nil, err
		}

		lo.log("graph generated")

		//copy generated graph and get name
		resString[2], err = copyGraph(config, condition.Itterator())
		if err != nil {
			return nil, err
		}

		lo.log("graph coppied")

		//start graph handler
		if err := startGraphHandler(config, condition.Itterator()); err != nil {
			return nil, err
		}

		lo.log("graph handler finnished")

		//save result and get name
		resString[3], err = saveGraphHandlerResult(config, condition.Itterator())
		if err != nil {
			return nil, err
		}

		lo.log("result coppied")

		resString[5], err = getResult(config)
		if err != nil {
			return nil, err
		}

		lo.log("result add to tab")

		//get string with time
		resString[4], err = getTime(config)
		if err != nil {
			return nil, err
		}

		//get time
		ittime, err = time.ParseDuration(resString[4])
		if err != nil {
			return nil, err
		}

		lo.log("time parsed")

		//get string with mark
		if config.MTCFG.WithFMMark {
			resString[6], err = getMark(config)
			if err != nil {
				return nil, err
			}
			lo.log("mark parsed")
		}

		//get string with advtime if enable
		if config.ATCFG.EnableAdvTime {
			advtime := ""
			if config.TypeOfTest == EDGETEST {
				advtime += strconv.Itoa(amE)
			}
			if config.TypeOfTest == VERTEXTEST {
				advtime += strconv.Itoa(amV)
			}
			tmp, err := getAdvancedTime(config)
			if err != nil {
				return nil, err
			}
			advtime += " " + tmp[:(len(tmp)-1)]
			lo.log("advtime parsed")

			if err = AppendStringToFile(PathToAdvTimeFile, advtime, condition.Itterator()); err != nil {
				return nil, err
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
			return nil, err
		}

		lo.log("result added to file")
	}
	return condition, nil
}
