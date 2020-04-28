package main

import (
	"encoding/csv"
	"os"
	"strings"

	csplotter "github.com/pplcc/plotext/custplotter"
)

func getCandels(config *ExtraGraphicCfg) (csplotter.TOHLCVer, error) {
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

	candels := make(csplotter.TOHLCVs, len(tab)-1)

	for i := 1; i < len(tab); i++ {
		operandT, err := parseOperand(m, config.NameFilds[0], config, tab[i])
		if err != nil {
			return nil, err
		}

		operandO, err := parseOperand(m, config.NameFilds[1], config, tab[i])
		if err != nil {
			return nil, err
		}

		operandH, err := parseOperand(m, config.NameFilds[2], config, tab[i])
		if err != nil {
			return nil, err
		}
		operandL, err := parseOperand(m, config.NameFilds[3], config, tab[i])
		if err != nil {
			return nil, err
		}
		operandC, err := parseOperand(m, config.NameFilds[4], config, tab[i])
		if err != nil {
			return nil, err
		}
		operandV, err := parseOperand(m, config.NameFilds[5], config, tab[i])
		if err != nil {
			return nil, err
		}

		candels[i-1].T = operandT
		candels[i-1].O = operandO
		candels[i-1].L = operandL
		candels[i-1].V = operandV

		nonzero := strings.Contains(config.Flag, "nonzero")
		if nonzero && operandH == 0.0 {
			candels[i-1].H = operandL
		} else {
			candels[i-1].H = operandH
		}

		if nonzero && operandH == 0.0 {
			candels[i-1].C = operandO
		} else {
			candels[i-1].C = operandC
		}

	}

	return candels, nil
}
