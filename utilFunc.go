package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type logger struct {
	state bool
}

func (l *logger) log(str string) {
	if l.state {
		log.Println(str)
	}
}

func createResDir(config *TestConfig) error {
	if config.SaveResultOfGraphHandlerFlag {
		if config.PathToFileWithResult == "" {
			return errors.New("wrong path to file with result")
		}
		if config.PathToDirForResult != "" {
			//creating custom dir if needed
			if _, err := os.Stat(config.PathToDirForResult); os.IsNotExist(err) {
				if err := os.MkdirAll(config.PathToDirForResult, os.ModePerm); err != nil {
					return err
				}
			}
		} else {
			//creating deffult dir
			if _, err := os.Stat("GraphResult"); os.IsNotExist(err) {
				if err := os.MkdirAll("GraphResult", os.ModePerm); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func createGraphDir(config *TestConfig) error {
	if config.SavingGeneratedGraphFlag {
		if config.GraphPath == "" {
			return errors.New("wrong path to generated graph")
		}
		if config.PathToDirForGeneratedGraph != "" {
			if _, err := os.Stat(config.PathToDirForGeneratedGraph); os.IsNotExist(err) {
				if err := os.MkdirAll(config.PathToDirForGeneratedGraph, os.ModePerm); err != nil {
					return err
				}
			}
		} else {
			if _, err := os.Stat("GeneratedGraphs"); os.IsNotExist(err) {
				if err := os.MkdirAll("GeneratedGraphs", os.ModePerm); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func countAmount(config *TestConfig, it int) (int, int) {
	return config.StartingAmountOfVertex + it*config.VertexDifferens, config.StartingAmountOfEdges + it*config.EdgesDifferens
}

func startGraphGene(config *TestConfig, it int) error {
	//count amount of vertex and edges on itteration
	amountOfVertex, amountOfEdges := countAmount(config, it)

	//generate configs
	flags := make([]string, len(config.GGCFG.GraphGeneratorFlags))
	copy(flags, config.GGCFG.GraphGeneratorFlags)
	for i, flag := range flags {
		if flag == "Avertex" {
			flags[i] = strconv.Itoa(amountOfVertex)
		}
		if flag == "Aedges" {
			flags[i] = strconv.Itoa(amountOfEdges)
		}
		if flag == "GraphPath" {
			flags[i] = config.GraphPath
		}
	}

	var graphgene *exec.Cmd
	if config.GGCFG.GraphGeneratorType != "" {
		ggFlags := ""
		for _, flag := range flags {
			ggFlags += flag + " "
		}

		stVMFlags := append(config.GGCFG.VMStarterFlags, config.GGCFG.GraphGeneratorType+" "+config.GGCFG.GraphGeneratorVMFlags+" "+config.GGCFG.GraphGeneratorPath+" "+ggFlags)

		graphgene = exec.Command(config.GGCFG.VMStarter, stVMFlags...)
	} else {
		graphgene = exec.Command(config.GGCFG.GraphGeneratorPath, flags...)
	}

	var stdout io.ReadCloser
	if config.GGCFG.Output {
		var err error
		stdout, err = graphgene.StdoutPipe()
		if err != nil {
			return err
		}
	}

	stderr, err := graphgene.StderrPipe()
	if err != nil {
		return err
	}

	if err := graphgene.Run(); err != nil {
		return err
	}

	errout, _ := ioutil.ReadAll(stderr)
	if string(errout) != "" {
		log.Printf("%s\n", errout)
	}

	if config.GGCFG.Output {
		out, _ := ioutil.ReadAll(stdout)
		if string(out) != "" {
			log.Printf("%s\n", out)
		}
	}

	return nil
}

//returning name of copy file and error or nil
func copyGraph(config *TestConfig, itterator int) (string, error) {
	//if need save
	if config.SavingGeneratedGraphFlag {
		if config.GraphPath != "" {
			if config.PathToDirForGeneratedGraph != "" {
				//save to custom dir
				if err := Copy(config.GraphPath, config.PathToDirForGeneratedGraph+"/graph_"+strconv.Itoa(itterator)); err != nil {
					return "", nil
				}
				return config.PathToDirForGeneratedGraph + "/graph_" + strconv.Itoa(itterator), nil
			} else {
				//save to def dir
				if err := Copy(config.GraphPath, "GeneratedGraphs"+"/graph_"+strconv.Itoa(itterator)); err != nil {
					return "", err
				}
				return "GeneratedGraphs" + "/graph_" + strconv.Itoa(itterator), nil
			}
		} else {
			return "-1", nil
		}
	} else {
		return "-1", nil
	}
}

func startGraphHandler(config *TestConfig, it int) error {
	//count amount of vertex and edges on itteration
	amountOfVertex, amountOfEdges := countAmount(config, it)

	//generate configs
	flags := make([]string, len(config.GHCFG.GraphHandlerFlags))
	copy(flags, config.GHCFG.GraphHandlerFlags)
	for i, flag := range flags {
		if flag == "Avertex" {
			flags[i] = strconv.Itoa(amountOfVertex)
		}
		if flag == "Aedges" {
			flags[i] = strconv.Itoa(amountOfEdges)
		}
		if flag == "GraphPath" {
			flags[i] = config.GraphPath
		}
	}

	var graphhandler *exec.Cmd
	if config.GHCFG.GraphHandlerType != "" {
		ggFlags := ""
		for _, flag := range flags {
			ggFlags += flag + " "
		}

		stVMFlags := append(config.GHCFG.VMStarterFlags, config.GHCFG.GraphHandlerType+" "+config.GHCFG.GraphHandlerVMFlags+" "+config.GHCFG.GraphHandlerPath+" "+ggFlags)

		graphhandler = exec.Command(config.GHCFG.VMStarter, stVMFlags...)
	} else {
		graphhandler = exec.Command(config.GHCFG.GraphHandlerPath, flags...)
	}

	var stdout io.ReadCloser
	if config.GHCFG.Output {
		var err error
		stdout, err = graphhandler.StdoutPipe()
		if err != nil {
			return err
		}
	}

	stderr, err := graphhandler.StderrPipe()
	if err != nil {
		return err
	}

	if err := graphhandler.Run(); err != nil {
		return err
	}

	errout, _ := ioutil.ReadAll(stderr)
	if string(errout) != "" {
		log.Printf("%s\n", errout)
	}

	if config.GHCFG.Output {
		out, _ := ioutil.ReadAll(stdout)
		if string(out) != "" {
			log.Printf("%s\n", out)
		}
	}

	return nil
}

func saveGraphHandlerResult(config *TestConfig, it int) (string, error) {
	if config.SaveResultOfGraphHandlerFlag {
		//save to cutom dir
		if config.PathToDirForResult != "" {
			if err := Copy(config.PathToFileWithResult, config.PathToDirForResult+"/ResultedGraph_"+strconv.Itoa(it)); err != nil {
				return "", err
			}
			return config.PathToDirForResult + "/ResultedGraph_" + strconv.Itoa(it), nil
		} else {
			//save to def dir
			if err := Copy(config.PathToFileWithResult, "GraphResult"+"/Resultedraph_"+strconv.Itoa(it)); err != nil {
				return "", err
			}
			return "GraphResult" + "/ResultedGraph_" + strconv.Itoa(it), nil
		}
	} else {
		return "-1", nil
	}
}

func getTime(config *TestConfig) (string, error) {
	//open file to read duration of gh work
	timeFile, err := os.Open(config.PathToFileWithTime)
	if err != nil {
		return "", err
	}
	defer timeFile.Close()

	//read duration
	scanner := bufio.NewScanner(timeFile)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", scanner.Err()
	} else {
		return scanner.Text(), nil
	}
}
