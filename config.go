package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type TestConfig struct {
	TypeOfTest string `yaml:"type_of_test"`

	OutputFlag bool `yaml:"output_flag"`

	GraphicTitle string `yaml:"graphic_title"`

	StartingAmountOfVertex int `yaml:"starting_amount_of_vertex"`
	VertexDifferens        int `yaml:"vertex_differens"`

	StartingAmountOfEdges int `yaml:"starting_amount_of_edges"`
	EdgesDifferens        int `yaml:"edges_differens"`

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

	GGCFG GraphGeneratorConfig `yaml:"graphgenerator_config"`
	GHCFG GraphHandlerConfig   `yaml:"graphhandler_config"`
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
