package main

import (
	"io"
	"log"
	"os"
)

var ResultNames = []string{"vertex", "edges", "path to graph", "path to result", "time", "result"}
var AdvTimeNames = []string{""}

func initResultNames(config *TestConfig) {
	if config.MTCFG.WithFMMark {
		ResultNames = append(ResultNames, "mark")
	}
	if err := resultwriter.Write(ResultNames); err != nil {
		log.Panic(err)
	}
}

func initAdvTimeNames(config *TestConfig) {
	if config.ATCFG.EnableAdvTime {
		switch config.TypeOfTest {
		case EDGETEST:
			AdvTimeNames[0] = "edges"
		case VERTEXTEST:
			AdvTimeNames[0] = "vertex"
		case ITTEST:
			AdvTimeNames[0] = "itteration"
		default:
			return
		}

		AdvTimeNames = append(AdvTimeNames, config.ATCFG.GraphicCFG.NameSet...)

		if err := advtimewriter.Write(AdvTimeNames); err != nil {
			log.Panic(err)
		}
	}
}

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

	//create file for results
	resFileCsv, err := os.Create(PathToResultFileCsv)
	if err != nil {
		resFile.Close()
		log.Panic(err)
	}
	resultwriter = NewCustWriterF(resFileCsv)

	lo.log("result file created")

	//create file for advtime
	advtimeFileCsv, err := os.Create(PathToAdvTimeFileCsv)
	if err != nil {
		advtimeFile.Close()
		log.Panic(err)
	}
	advtimewriter = NewCustWriterF(advtimeFileCsv)
}

func initWriters() {
	advtimewriter = NewCustWriter(PathToAdvTimeFileCsv)
	resultwriter = NewCustWriter(PathToResultFileCsv)
}

func drawGraphics(config *TestConfig, condition *TestState) error {
	//open file with results
	resFile, err := os.Open(PathToResultFile)
	if err != nil {
		return err
	}
	defer resFile.Close()

	//draw graphic
	err = DrawPlotCust(resFile, config, condition.Itterator())
	if err != nil {
		return err
	}

	if config.MTCFG.DrawDiffGraphic {
		lo.log("draw diff graphic")
		resFile.Seek(0, io.SeekStart)
		err = DrawMarkDiff(resFile, config, condition.Itterator())
		if err != nil {
			return err
		}
	}
	if config.MTCFG.DrawDynGraphic {
		resFile.Seek(0, io.SeekStart)
		lo.log("draw progression graphic")
		err = DrawMarkProgression(resFile, config, condition.Itterator())
		if err != nil {
			return err
		}
	}

	if config.ATCFG.DrawDistribGraphic {
		advtimeFile, err := os.Open(PathToAdvTimeFile)
		if err != nil {
			return err
		}
		defer advtimeFile.Close()
		lo.log("draw distribution graphic")
		err = DrawAdvtimeDistribution(advtimeFile, config, condition.Itterator())
		if err != nil {
			return err
		}
	}

	return nil
}
