package main

import (
	"log"
	"os"
	"strconv"
	"io"
)

func main() {
	config, err := readConfig(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	file, err := os.Open("ResultTab")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	am, _ := strconv.Atoi(os.Args[2])
	if err = DrawPlotCust(file, config, am); err != nil {
		log.Println(err)
	}

	if config.MTCFG.DrawDiffGraphic{
		file.Seek(0, io.SeekStart)
		log.Println("draw diff graphic")
		err = DrawMarkDiff(file, config, am)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if config.MTCFG.DrawDynGraphic{
		file.Seek(0, io.SeekStart)
		log.Println("draw progression graphic")
		err = DrawMarkProgression(file, config, am)
		if err != nil {
			log.Println(err)
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
		log.Println("draw distribution graphic")
		err = DrawAdvtimeDistribution(advtimeFile,config,am)
		if err != nil{
			log.Println(err)
			return
		}
	}
}
