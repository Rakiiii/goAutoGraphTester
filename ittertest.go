package main

import (
	"strconv"
	"strings"
	"time"
)

func startItterationTest(config *TestConfig) (*TestState, error) {

	condition := NewTestState(config)
	ittime, _ := time.ParseDuration("0ms")
	size := *NewGraphSize()
	graphgeneFlag := true
	var err error

	for condition.isContinue(ittime, size.AE, size.AV) {
		resString := make([]string, 5)

		resString[0] = strconv.Itoa(config.ITCFG.StartingAmountOfItteration + config.ITCFG.ItterrationDifference*condition.Itterator())

		if config.ITCFG.GraphGeneratorBefavor != "no" {
			if graphgeneFlag {
				//generate graph
				if size, err = startGraphGene(config, condition.Itterator()); err != nil {
					return nil, err
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

		lo.log("result add to tab")

		//get result value
		resString[1], err = getResult(config)
		if err != nil {
			return nil, err
		}

		lo.log("result value getted")

		//get string with mark
		if config.MTCFG.WithFMMark {
			resString[5], err = getMark(config)
			if err != nil {
				return nil, err
			}
			lo.log("mark parsed")
		}

		//get string with advtime if enable
		if config.ATCFG.EnableAdvTime {
			advtime := strconv.Itoa(config.ITCFG.StartingAmountOfItteration + config.ITCFG.ItterrationDifference*condition.Itterator())
			tmp, err := getAdvancedTime(config)
			if err != nil {
				return nil, err
			}
			advtime += " " + tmp[:(len(tmp)-1)]
			lo.log("advtime parsed")

			if err = advtimewriter.Write(strings.Fields(advtime)); err != nil {
				return nil, err
			}

			/*if err = AppendStringToFile(PathToAdvTimeFile, advtime, condition.Itterator()); err != nil {
				return nil, err
			}*/
			lo.log("advtime appended")
		}

		if err := resultwriter.Write(resString); err != nil {
			return nil, err
		}

		/*
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
		*/
	}

	return condition, nil
}
