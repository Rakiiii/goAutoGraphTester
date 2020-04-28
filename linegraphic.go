package main

import (
	"encoding/csv"
	"os"
	"strings"

	"gonum.org/v1/plot/plotter"
)

func getXYSPoints(config *ExtraGraphicCfg) (plotter.XYer, error) {
	file, err := os.Open(config.PathToSoures)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tab, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	points := make(plotter.XYZs, len(tab)-1)

	m := createParamMap(tab[0])

	for i := 1; i < len(tab); i++ {
		points[i-1].X, err = parseOperand(m, config.NameFilds[0], config, tab[i])
		if err != nil {
			return nil, err
		}
		points[i-1].Y, err = parseOperand(m, config.NameFilds[1], config, tab[i])
		if err != nil {
			return nil, err
		}

		if strings.Contains(config.Flag, "positive") {
			if points[i-1].X < 0 {
				points[i-1].X = 0.0
			}
			if points[i-1].Y < 0 {
				points[i-1].Y = 0.0
			}
		}
	}

	return points, nil
}
