package main

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

const (
	PathToResultFile     = "ResultTab"
	PathToAdvTimeFile    = "AdvTimeTab"
	PathToResultFileCsv  = "ResultTab.csv"
	PathToAdvTimeFileCsv = "AdvTimeTab.csv"
	EDGETEST             = "edgestest"
	VERTEXTEST           = "vertextest"
	ITTEST               = "ittest"
	PARSETEST            = "parsetest"
	CMPTEST              = "cmptest"
	TIMESTOP             = "timestop"
	ITSTOP               = "itstop"
	EDGSTOP              = "edgstop"
	VERTEXSTOP           = "vertexstop"
	MIXEDSTOP            = "mixed"
	VERTEXFLAG           = "Avertex"
	EDGEFLAG             = "Aedges"
	GRAPHPATHFLAG        = "GraphPath"
	ITFLAG               = "it"
	PARSEGRAPHFLAG       = "pgraph"
)

type TestConfig struct {
	TypeOfTest string `yaml:"type_of_test"`

	OutputFlag bool `yaml:"output_flag"`

	GraphicTitle string `yaml:"graphic_title"`

	StartingAmountOfVertex int `yaml:"starting_amount_of_vertex"`
	VertexDifferens        int `yaml:"vertex_differens"`
	MaxAmountOfVertex      int `yaml:max_amount_of_vertex`

	StartingAmountOfEdges int  `yaml:"starting_amount_of_edges"`
	EdgesDifferens        int  `yaml:"edges_differens"`
	RndEdges              bool `yaml:rnd_edges`
	MaxAmountOfEdges      int  `yaml:max_amount_of_edges`

	TypeOfStopCondition  string `yaml:"type_of_stop_condition"`
	MaxTimeForItteration int64  `yaml:"max_time_for_itteration"`
	AmountOfItterations  int    `yaml:"amount_of_itterations"`

	PathToFileWithTime string `yaml:"path_to_file_with_time"`

	SaveResultOfGraphHandlerFlag bool   `yaml:"save_result_flag"`
	PathToDirForResult           string `yaml:"path_to_dir_for_result"`
	PathToFileWithResult         string `yaml:"path_to_file_with_result"`

	SavingGeneratedGraphFlag   bool   `yaml:"save_generated_graph_flag"`
	PathToDirForGeneratedGraph string `yaml:"path_to_dir_for_coping_generated_graph"`
	GraphPath                  string `yaml:"graphpath"`

	GGCFG GraphGeneratorConfig  `yaml:"graphgenerator_config"`
	GHCFG GraphHandlerConfig    `yaml:"graphhandler_config"`
	ITCFG ItterrationTestConfig `yaml:"ittest_config"`
	PTCFG ParseTestConfig       `yaml:"parsetest_config"`
	MTCFG MarkTestConfig        `yaml:"marktest_config"`
	ATCFG AdvTimeConfig         `yaml:"advtime_config"`

	GraphicSet []ExtraGraphicCfg
}

type GraphGeneratorConfig struct {
	GraphGeneratorPath    string   `yaml:"graphgenerator_path"`
	GraphGeneratorFlags   []string `yaml:"graphgenerator_flags"`
	GraphGeneratorType    string   `yaml:"graphfenerator_type"`
	GraphGeneratorVMFlags string   `yaml:"graphgenerator_vm_flags"`
	VMStarter             string   `yaml:"vmstartergg"`
	VMStarterFlags        []string `yaml:"vmstartergg_flags"`
	Output                bool     `yaml:"ggoutput"`
}

type GraphHandlerConfig struct {
	GraphHandlerPath    string   `yaml:"graphhandler_path"`
	GraphHandlerFlags   []string `yaml:"graphhandler_flags"`
	GraphHandlerType    string   `yaml:"graphhandler_type"`
	GraphHandlerVMFlags string   `yaml:"graphhandler_vm_flags"`
	VMStarter           string   `yaml:"vmstartergh"`
	VMStarterFlags      []string `yaml:"vmstartergh_flags"`
	Output              bool     `yaml:"ghoutput"`
}

type ItterrationTestConfig struct {
	GraphGeneratorBefavor      string `yaml:"graphgenerator_behavor"`
	PathToFileWithResult       string `yaml:"result_path"`
	StartingAmountOfItteration int    `yaml:"start_amount_of_itteration"`
	ItterrationDifference      int    `yaml:"itteration_difference"`
}

type ParseTestConfig struct {
	PathToDir            string `yaml:"path_to_dir_with_files"`
	FileMask             string `yaml:"file_mask"`
	PathToFileWithResult string `yaml:"result_path_parsed"`
	File                 string
}

type MarkTestConfig struct {
	WithFMMark      bool   `yaml:"contains_mark"`
	PathToFile      string `yaml:"path_to_file"`
	DrawDiffGraphic bool   `yaml:"draw_diff_graphic"`
	DrawDynGraphic  bool   `yaml:"draw_dyn_graphic"`
}

type AdvTimeConfig struct {
	EnableAdvTime      bool             `yaml:"enable_adv_time"`
	PathToFile         string           `yaml:"path_to_file"`
	DrawDistribGraphic bool             `yaml:"draw_distribution_graphic"`
	GraphicCFG         AdvGraphicConfig `yaml:"adv_graphic_config"`
}

type AdvGraphicConfig struct {
	ColorSet []string `yaml:"color_set"`
	NameSet  []string `yaml:"name_set"`
}

type ExtraGraphicCfg struct {
	//name of file with result graphic
	Name string
	//must be csv file
	PathToSoures string
	//name of fileds from csv which will be used for graphic draw
	//for "line" first one is X cord second is y cord
	//for "candels" is represented by TOHLCV  [ first is T, second is O and so on]
	//"multybars" is represented by set of bars stacked on each other [ firstone is x axis value,second is first part of bar,third is second part of bar and so on]
	NameFilds []string
	//Extra operation, must looks like [[operation1 nameOfOperand1 nameOfOperand2],[operation2 nameOfOperand1 nameOfOperand2]] to use result in graphic
	//must right "!oper1" to NameFields,
	//Possibel operations sub,plus,div,time
	Operation []string
	//Type of graphic chich will be drawen:"line","candels","multybars"
	//for "multybars" CFG must be non nil
	Type string
	//label of x asix
	XAsixLabel string
	//label of y asix
	YAsixLabel string
	//label of z asix
	GraphicLabel string
	//must contains name for legend and color set
	//for "candels" first color is upper color, second is lower color
	CFG AdvGraphicConfig
	//string that contains 2 words with space : first : "top" or "bottom" second : "left" or "right"
	//default top,left
	LegendPosition string
	//true if need to add legend
	DoLegend bool

	//some extra falgs, must be separate with spaces
	//if "nonzero" setted and candels flag than if C or H is zero then they will be setted the same way as O and L
	//if "positive" setted and line then x or y cord negative values would be setted as 0.0
	//if "inv" setted and candels flag than colors moved between upped and lowwer
	//if "nonzero" setted and "multybars" than must look "nonzero=[some number]" it will be YMin
	//if "length" setted and "multybars" than must look like "length=[some number]"[it will be vg.Length()]
	Flag string
}

func readConfig(configName string) (*TestConfig, error) {
	file, err := ioutil.ReadFile(configName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	config := new(TestConfig)

	if err = yaml.Unmarshal(file, config); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return config, nil
}
