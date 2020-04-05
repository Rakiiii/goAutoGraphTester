package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"io"
	//"bufio"
)

//var PathToResultFile string = "ResultTab"

var lo logger = logger{state:true}

func main() {
	//reading config
	config, err := readConfig(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	/*fmt.Println("enable adv time:",config.ATCFG.EnableAdvTime)
	fmt.Println("enable draw graphic:",config.ATCFG.DrawDistribGraphic)
	fmt.Print("Press 'Enter' to continue...")
  	bufio.NewReader(os.Stdin).ReadBytes('\n')*/

		

	//create dir for result of graph handler
	if err := createResDir(config); err != nil {
		log.Println(err)
		return
	}

	lo.log("res dir created")

	//create dir for generated graphs
	if err := createGraphDir(config); err != nil {
		log.Println(err)
		return
	}

	lo.log("graph dir created")

	//create file for results
	resFile, err := os.Create(PathToResultFile)
	if err != nil {
		fmt.Println(err)
		resFile.Close()
		return
	}
	resFile.Close()

	lo.log("result file created")

	//create file for results
	advtimeFile, err := os.Create(PathToAdvTimeFile)
	if err != nil {
		fmt.Println(err)
		advtimeFile.Close()
		return
	}
	advtimeFile.Close()

	var prevTime time.Duration
	var itterator int = 0
	timeFlag := true
	itterationFlag := true
	graphgeneFlag := true

	for timeFlag && itterationFlag {
		switch {
		case config.TypeOfTest == EDGETEST || config.TypeOfTest == VERTEXTEST:

			lo.log("start vertex of edges test")

			resString := make([]string, 7)

			amV, amE := countAmount(config, itterator)
			resString[0] = strconv.Itoa(amV)
			resString[1] = strconv.Itoa(amE)

			log.Println("start graph generation")
			//generate graph
			if err := startGraphGene(config, itterator); err != nil {
				log.Println(err)
				return
			}

			lo.log("graph generated")

			//copy generated graph and get name
			resString[2], err = copyGraph(config, itterator)
			if err != nil {
				log.Println(err)
				return
			}

			lo.log("graph coppied")

			//start graph handler
			if err := startGraphHandler(config, itterator); err != nil {
				log.Println(err)
				return
			}

			lo.log("graph handler finnished")

			//save result and get name
			resString[3], err = saveGraphHandlerResult(config, itterator)
			if err != nil {
				log.Println(err)
				return
			}

			lo.log("result coppied")

			resString[5],err = getResult(config)
			if err != nil {
				log.Println(err)
				return
			}

			lo.log("result add to tab")

			//get string with time
			resString[4], err = getTime(config)
			if err != nil {
				log.Println(err)
				return
			}

			//get time
			prevTime, err = time.ParseDuration(resString[4])
			if err != nil {
				log.Println(err)
				return
			}

			lo.log("time parsed")


			//get string with mark
			if config.MTCFG.WithFMMark{
				resString[6], err = getMark(config)
				if err != nil {
					log.Println(err)
					return
				}	
				lo.log("mark parsed")
			}

			//get string with advtime if enable
			if config.ATCFG.EnableAdvTime{
				advtime := ""
				if config.TypeOfTest == EDGETEST {
					advtime += strconv.Itoa(amE)
				}
				if config.TypeOfTest == VERTEXTEST{
					advtime += strconv.Itoa(amV)
				}
				tmp,err := getAdvancedTime(config)
				if err != nil{
					log.Println(err)
					return
				}
				advtime += " " + tmp[:(len(tmp)-1)]
				lo.log("advtime parsed")

				if err = AppendStringToFile(PathToAdvTimeFile,advtime,itterator);err != nil{
					log.Println(err)
					return
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
			if err = AppendStringToFile(PathToResultFile, writeRes, itterator); err != nil {
				log.Println(err)
				return
			}		

			lo.log("result added to file")



		case config.TypeOfTest == ITTEST:
			resString := make([]string, 5)

			resString[0] = strconv.Itoa(config.ITCFG.StartingAmountOfItteration + config.ITCFG.ItterrationDifference*itterator)

			if config.ITCFG.GraphGeneratorBefavor != "no" {
				if graphgeneFlag {
					//generate graph
					if err := startGraphGene(config, itterator); err != nil {
						log.Println(err)
						return
					}
					lo.log("graph generated")
				}
			}

			if config.ITCFG.GraphGeneratorBefavor == "once" {
				graphgeneFlag = false
			}

			//copy generated graph and get name
			resString[2], err = copyGraph(config, itterator)
			if err != nil {
				log.Println(err)
				return
			}

			lo.log("graph coppied")

			//start graph handler
			if err := startGraphHandler(config, itterator); err != nil {
				log.Println(err)
				return
			}

			lo.log("graph handler finnished")

			//save result and get name
			resString[3], err = saveGraphHandlerResult(config, itterator)
			if err != nil {
				log.Println(err)
				return
			}
			

			lo.log("result coppied")

			lo.log("result add to tab")

			//get result value
			resString[1], err = getResult(config)
			if err != nil {
				log.Println(err)
				return
			}

			lo.log("result value getted")

			//get string with mark
			if config.MTCFG.WithFMMark{
				resString[5], err = getMark(config)
				if err != nil {
					log.Println(err)
					return
				}	
				lo.log("mark parsed")
			}

						//get string with advtime if enable
			if config.ATCFG.EnableAdvTime{
				advtime := strconv.Itoa(config.ITCFG.StartingAmountOfItteration + config.ITCFG.ItterrationDifference*itterator)
				tmp,err := getAdvancedTime(config)
				if err != nil{
					log.Println(err)
					return
				}
				advtime +=  " " +  tmp[:(len(tmp)-1)]
				lo.log("advtime parsed")
				if err = AppendStringToFile(PathToAdvTimeFile,advtime,itterator);err != nil{
					log.Println(err)
					return
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
			if err = AppendStringToFile(PathToResultFile, writeRes, itterator); err != nil {
				log.Println(err)
				return
			}
		case config.TypeOfTest == PARSETEST:

			resString := make([]string, 2)

			setOfGraphs, err := getSliceOfGrpahs(config)
			if err != nil {
				log.Panicln(err)
				return
			}

			log.Println("set of grpahs parsed")

			for it, graphFile := range *setOfGraphs {
				log.Println("grpah handling starts , grpah name:", graphFile)
				config.PTCFG.File = graphFile
				if err := startGraphHandler(config, it); err != nil {
					log.Panicln(err)
					return
				}

				lo.log("graph handler finnished")

				//save result and get name
				resString[0], err = saveGraphHandlerResult(config, itterator)
				if err != nil {
					log.Println(err)
					return
				}

				lo.log("result coppied")

				//get result value
				resString[1], err = getResultFromParsed(config)
				if err != nil {
					log.Println(err)
					return
				}
				lo.log("result value getted")

				//get string with mark
				resString[5], err = getMark(config)
				if err != nil {
					log.Println(err)
					return
				}	

				if config.MTCFG.WithFMMark{
					lo.log("mark parsed")
				}

				//make res string
				writeRes := ""
				for _, i := range resString {
					writeRes += i
					writeRes += " "
				}

				//write result string
				if err = AppendStringToFile(PathToResultFile, writeRes, itterator); err != nil {
					log.Println(err)
					return
				}

			}

		default:
			log.Panicln("Wrong type of test")
			return
		}

		//stop condition
		if (config.TypeOfStopCondition == TIMESTOP || config.TypeOfStopCondition == MIXEDSTOP) && prevTime.Milliseconds() > config.MaxTimeForItteration {
			timeFlag = false
		}

		if (config.TypeOfStopCondition == ITSTOP || config.TypeOfStopCondition == MIXEDSTOP) && config.AmountOfItterations <= itterator {
			itterationFlag = false
		}

		itterator++
	}

	if config.TypeOfTest != PARSETEST {
		//open file with results
		resFile, err = os.Open(PathToResultFile)
		if err != nil {
			log.Println(err)
			return
		}
		defer resFile.Close()

		//draw graphic
		err = DrawPlotCust(resFile, config, itterator)
		if err != nil {
			fmt.Println(err)
			return
		}

		if config.MTCFG.DrawDiffGraphic{
			lo.log("draw diff graphic")
			resFile.Seek(0, io.SeekStart)
			err = DrawMarkDiff(resFile, config, itterator)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if config.MTCFG.DrawDynGraphic{
			resFile.Seek(0, io.SeekStart)
			lo.log("draw progression graphic")
			err = DrawMarkProgression(resFile, config, itterator)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if config.ATCFG.DrawDistribGraphic {
			advtimeFile,err := os.Open(PathToAdvTimeFile)
			if err != nil{
				log.Println(err)
				return
			}
			defer advtimeFile.Close()
			lo.log("draw distribution graphic")
			err = DrawAdvtimeDistribution(advtimeFile,config,itterator)
			if err != nil{
				log.Println(err)
				return
			}
		}
	}
}
