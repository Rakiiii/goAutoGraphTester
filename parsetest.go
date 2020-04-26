package main

import (
	"log"
	"time"
)

func startParseTest(config *TestConfig) (*TestState, error) {

	ittime, _ := time.ParseDuration("0ms")

	resString := make([]string, 2)

	setOfGraphs, err := getSliceOfGrpahs(config)
	if err != nil {
		return nil, err
	}

	condition := NewTestState(config)

	log.Println("set of grpahs parsed")

	for it, graphFile := range *setOfGraphs {
		condition.isContinue(ittime, 0, 0)

		log.Println("grpah handling starts , grpah name:", graphFile)
		config.PTCFG.File = graphFile
		if err := startGraphHandler(config, it); err != nil {
			return nil, err
		}

		lo.log("graph handler finnished")

		//save result and get name
		resString[0], err = saveGraphHandlerResult(config, condition.Itterator())
		if err != nil {
			return nil, err
		}

		lo.log("result coppied")

		//get result value
		resString[1], err = getResultFromParsed(config)
		if err != nil {
			return nil, err
		}
		lo.log("result value getted")

		//get string with mark
		resString[5], err = getMark(config)
		if err != nil {
			return nil, err
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
			return nil, err
		}

	}

	return condition, nil
}
