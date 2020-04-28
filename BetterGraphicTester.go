package main

import (
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
	initAdvTimeNames(config)
	initResultNames(config)
	config.GraphicSet = initStdGraphics(config)

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

	flushWriters()

	if config.TypeOfTest != PARSETEST {
		if err = drawGraphics(config, condition); err != nil {
			log.Panic(err)
		}
	}
}
