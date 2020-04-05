package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	DEFFRESDIR = "GraphResult"
	DEFFGRPAHDIR = "GeneratedGraphs"
	DEFFGRAPHNAME = "/graph_"
	DEFFRESNAME = "/ResultedGraph_"
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
			if _, err := os.Stat(DEFFRESDIR); os.IsNotExist(err) {
				if err := os.MkdirAll(DEFFRESDIR, os.ModePerm); err != nil {
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
			if _, err := os.Stat(DEFFGRPAHDIR); os.IsNotExist(err) {
				if err := os.MkdirAll(DEFFGRPAHDIR, os.ModePerm); err != nil {
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
	flags := *parseFlags(&config.GGCFG.GraphGeneratorFlags, amountOfVertex, amountOfEdges, it, config, config.PTCFG.File)

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
		log.Println(flags)
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
				if err := Copy(config.GraphPath, config.PathToDirForGeneratedGraph+DEFFGRAPHNAME+strconv.Itoa(itterator)); err != nil {
					return "", nil
				}
				return config.PathToDirForGeneratedGraph + DEFFGRAPHNAME + strconv.Itoa(itterator), nil
			} else {
				//save to def dir
				if err := Copy(config.GraphPath, DEFFGRPAHDIR+DEFFGRAPHNAME+strconv.Itoa(itterator)); err != nil {
					return "", err
				}
				return DEFFGRPAHDIR + DEFFGRAPHNAME + strconv.Itoa(itterator), nil
			}
		} else {
			return "-1", nil
		}
	} else {
		return "-1", nil
	}
}

func startGraphHandler(config *TestConfig, it int) error {
	var lo logger = logger{config.OutputFlag}
	//count amount of vertex and edges on itteration
	amountOfVertex, amountOfEdges := countAmount(config, it)

	//generate configs
	flags := *parseFlags(&config.GHCFG.GraphHandlerFlags, amountOfVertex, amountOfEdges, it, config, config.PTCFG.File)

	fmt.Println(flags)

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

	lo.log("std out start")

	//var stdout io.ReadCloser
	if config.GHCFG.Output {

		graphhandler.Stdout = os.Stdout

	}

	lo.log("stderr start")
	graphhandler.Stderr = os.Stderr


	lo.log("handler started with out vm")
	if err := graphhandler.Run(); err != nil {
		return err
	}
	lo.log("handler finished")

	lo.log("stderr end")

	lo.log("stdout stoped")

	return nil
}

func saveGraphHandlerResult(config *TestConfig, it int) (string, error) {
	if config.SaveResultOfGraphHandlerFlag {
		//save to cutom dir
		if config.PathToDirForResult != "" {
			if err := Copy(config.PathToFileWithResult, config.PathToDirForResult+DEFFRESNAME+strconv.Itoa(it)); err != nil {
				return "", err
			}
			return config.PathToDirForResult + DEFFRESNAME + strconv.Itoa(it), nil
		} else {
			//save to def dir
			if err := Copy(config.PathToFileWithResult, DEFFRESDIR+DEFFRESNAME+strconv.Itoa(it)); err != nil {
				return "", err
			}
			return DEFFRESDIR + DEFFRESNAME + strconv.Itoa(it), nil
		}
	} else {
		return "-1", nil
	}
}

func getResult(config *TestConfig) (string, error) {
	if config.ITCFG.PathToFileWithResult == "" {
		file, err := os.Open(config.PathToFileWithResult)
		if err != nil {
			return "", err
		}
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		return scanner.Text(), nil

	} else {
		file, err := os.Open(config.ITCFG.PathToFileWithResult)
		if err != nil {
			return "", err
		}
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		return scanner.Text(), nil
	}
}

func getResultFromParsed(config *TestConfig) (string, error) {
	if config.PTCFG.PathToFileWithResult == "" {
		file, err := os.Open(config.PathToFileWithResult)
		if err != nil {
			return "", err
		}
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		return scanner.Text(), nil

	} else {
		file, err := os.Open(config.PTCFG.PathToFileWithResult)
		if err != nil {
			return "", err
		}
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		return scanner.Text(), nil
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

func parseFlags(pars *[]string, amountOfVertex, amountOfEdges, it int, config *TestConfig, parsedgraph string) *[]string {
	flags := make([]string, len(*pars))
	copy(flags, *pars)
	for i, flag := range flags {
		if flag == VERTEXFLAG {
			flags[i] = strconv.Itoa(amountOfVertex)
		}
		if flag == EDGEFLAG {
			flags[i] = strconv.Itoa(amountOfEdges)
		}
		if flag == GRAPHPATHFLAG {
			flags[i] = config.GraphPath
		}
		if flag == ITFLAG{
			flags[i] = strconv.Itoa(config.ITCFG.StartingAmountOfItteration + config.ITCFG.ItterrationDifference*it)
		}
		if flag == PARSEGRAPHFLAG {
			flags[i] = parsedgraph
		}
	}
	return &flags
}

func getSliceOfGrpahs(config *TestConfig) (*[]string, error) {
	res, err := filepath.Glob(config.PTCFG.PathToDir + config.PTCFG.FileMask)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func getMark(config *TestConfig)(string,error){
	if config.MTCFG.WithFMMark{
		//open file to read duration of gh work
		markFile, err := os.Open(config.MTCFG.PathToFile)
		if err != nil {
			return "", err
		}
		defer markFile.Close()
	
		//read duration
		scanner := bufio.NewScanner(markFile)
		scanner.Scan()
		if scanner.Err() != nil {
			return "", scanner.Err()
		} else {
			return scanner.Text(), nil
		}
	}else{
		return "",nil
	}
}

func getAdvancedTime(config *TestConfig)(string,error){
	if config.ATCFG.EnableAdvTime{
		advtimeFile,err := os.Open(config.ATCFG.PathToFile)
		if err != nil{
			return "",err
		}
		defer advtimeFile.Close()
		
		res := ""
		scanner := bufio.NewScanner(advtimeFile)
		for scanner.Scan(){
			res += scanner.Text() + "\n"
		}
		return res,nil
	}else{
		return "",nil
	}
}
