package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	mbplotter "github.com/Rakiiii/goMultiBarPlotter"
)

var WrongFlagMB = errors.New("Wrong flag in multy bar plotter")

func getMultyBars(config *ExtraGraphicCfg) (mbplotter.Barer, error) {
	file, err := os.Open(config.PathToSoures)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tab, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	m := createParamMap(tab[0])

	bars := make(mbplotter.Bars, len(tab)-1)
	var ymin float64

	if strings.Contains(config.Flag, "nonzero") {
		ymin, err = parseNumberFromFlag(config.Flag, "nonzero")
		if err != nil {
			return nil, err
		}
	} else {
		ymin = 0.0
	}

	//init slice of operands and get all operands from set
	for i := 1; i < len(tab); i++ {
		operandX, err := parseOperand(m, config.NameFilds[0], config, tab[i])
		if err != nil {
			return nil, err
		}

		operandSlice := make([]float64, len(config.NameFilds)-1)
		for j := 1; j < len(config.NameFilds); j++ {
			operandSlice[j-1], err = parseOperand(m, config.NameFilds[j], config, tab[i])
			if err != nil {
				return nil, err
			}
		}

		bars[i-1].Ymin = ymin
		bars[i-1].X = operandX

		bars[i-1].Y = operandSlice
	}

	return bars, nil
}
