package main

import (
	"log"
	"os"
)

var ResultNames = []string{"vertex", "edges", "path to graph", "path to result", "time", "result"}
var AdvTimeNames = []string{""}

func initResultNames(config *TestConfig) {
	if config.MTCFG.WithFMMark {
		ResultNames = append(ResultNames, "mark")
	}
	if config.TypeOfTest == ITTEST {
		ResultNames[1] = "itteration"
		ResultNames = ResultNames[1:]
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
	/*
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
	*/
	//create file for results
	resFileCsv, err := os.Create(PathToResultFileCsv)
	if err != nil {
		resFileCsv.Close()
		log.Panic(err)
	}
	resultwriter = NewCustWriterF(resFileCsv)

	lo.log("result file created")

	//create file for advtime
	advtimeFileCsv, err := os.Create(PathToAdvTimeFileCsv)
	if err != nil {
		advtimeFileCsv.Close()
		log.Panic(err)
	}
	advtimewriter = NewCustWriterF(advtimeFileCsv)
}

func initWriters() {
	advtimewriter = NewCustWriter(PathToAdvTimeFileCsv)
	resultwriter = NewCustWriter(PathToResultFileCsv)
}

func drawGraphics(config *TestConfig, condition *TestState) error {
	cfg := initMainGraphic(config)
	if err := drawGraphic(cfg); err != nil {
		return nil
	}

	cndcfg := initCandelGraphic(config)
	if err := drawGraphic(cndcfg); err != nil {
		return nil
	}

	dercfg := initMarkProgressionGraphic(config)
	if err := drawGraphic(dercfg); err != nil {
		return nil
	}

	advcfg := initAdvTimeGraphic(config)
	if err := drawGraphic(advcfg); err != nil {
		return err
	}

	/*
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
		}*/

	return nil
}

func initMainGraphic(config *TestConfig) *ExtraGraphicCfg {
	var lb string
	var x string
	switch config.TypeOfTest {
	case VERTEXTEST:
		lb = VERTEXTESTASICLABEL
		x = "vertex"
	case EDGETEST:
		lb = EDGETESTASICLABEL
		x = "edges"
	case ITTEST:
		lb = ITTESTASICLABEL
		x = "itteration"
	}

	return &ExtraGraphicCfg{Name: "MainGraphic.png",
		PathToSoures:   PathToResultFileCsv,
		NameFilds:      []string{x, "time"},
		Operation:      []string{},
		Type:           "line",
		XAsixLabel:     lb,
		YAsixLabel:     TIMEASICLABEL,
		GraphicLabel:   config.GraphicTitle,
		DoLegend:       false,
		LegendPosition: ""}
}

func initCandelGraphic(config *TestConfig) *ExtraGraphicCfg {
	var lb string
	var x string
	switch config.TypeOfTest {
	case VERTEXTEST:
		lb = VERTEXTESTASICLABEL
		x = "vertex"
	case EDGETEST:
		lb = EDGETESTASICLABEL
		x = "edges"
	case ITTEST:
		lb = ITTESTASICLABEL
		x = "itteration"
	}
	return &ExtraGraphicCfg{Name: "TestCandels.png",
		PathToSoures:   PathToResultFileCsv,
		NameFilds:      []string{x, "result", "mark", "result", "mark", "0.0"},
		Operation:      []string{},
		Type:           "candels",
		XAsixLabel:     lb,
		YAsixLabel:     "mark deviation from result",
		GraphicLabel:   "Deviation visualisation",
		DoLegend:       false,
		LegendPosition: "",
		Flag:           "inv nonzero",
	}
}

func initMarkProgressionGraphic(config *TestConfig) *ExtraGraphicCfg {
	var lb string
	var x string
	switch config.TypeOfTest {
	case VERTEXTEST:
		lb = VERTEXTESTASICLABEL
		x = "vertex"
	case EDGETEST:
		lb = EDGETESTASICLABEL
		x = "edges"
	case ITTEST:
		lb = ITTESTASICLABEL
		x = "itteration"
	}

	return &ExtraGraphicCfg{Name: "DerevationGraphic.png",
		PathToSoures:   PathToResultFileCsv,
		NameFilds:      []string{x, "!oper0"},
		Operation:      []string{"sub mark result"},
		Type:           "line",
		XAsixLabel:     lb,
		YAsixLabel:     "derivation",
		GraphicLabel:   "Derivative of mark",
		DoLegend:       false,
		LegendPosition: "",
		Flag:           "positive"}
}

func initAdvTimeGraphic(config *TestConfig) *ExtraGraphicCfg {
	var lb string
	var x string
	switch config.TypeOfTest {
	case VERTEXTEST:
		lb = VERTEXTESTASICLABEL
		x = "vertex"
	case EDGETEST:
		lb = EDGETESTASICLABEL
		x = "edges"
	case ITTEST:
		lb = ITTESTASICLABEL
		x = "itteration"
	}

	return &ExtraGraphicCfg{Name: "AdvTimeGraphic.png",
		PathToSoures:   PathToAdvTimeFileCsv,
		NameFilds:      append([]string{x}, config.ATCFG.GraphicCFG.NameSet...),
		Operation:      []string{},
		Type:           "multybars",
		XAsixLabel:     lb,
		YAsixLabel:     TIMEASICLABEL,
		GraphicLabel:   "Advanced time distribution",
		DoLegend:       true,
		LegendPosition: "top left",
		Flag:           "length=10",
		CFG:            config.ATCFG.GraphicCFG}
}
