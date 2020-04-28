package main

import (
	"errors"
	"image/color"
	"strconv"
	"strings"

	mbplotter "github.com/Rakiiii/goMultiBarPlotter"
	csplotter "github.com/pplcc/plotext/custplotter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var UnknownOperation = errors.New("Unknown operation for drawing graphic")

const (
	LINEPLOT   = "line"
	CANDELPLOT = "candels"
	MULTYBARS  = "multybars"

	SUB  = "sub"
	PLUS = "plus"
	DIV  = "div"
	TIME = "time"

	TOPPOS    = "top"
	BOTTOMPOS = "bottom"
	LEFTPOS   = "left"
	RIGHTPOS  = "right"
)

func drawGraphic(config *ExtraGraphicCfg) error {

	plotter, err := initPlotter(config)
	if err != nil {
		return err
	}
	plot, err := initPlot(config)
	if err != nil {
		return nil
	}

	plot.Add(plotter)

	switch plotter.(type) {
	case *mbplotter.MultiBarPlotter:
		plt := plotter.(*mbplotter.MultiBarPlotter)
		if config.DoLegend {
			for i := range plt.Colors {
				plot.Legend.Add(config.CFG.NameSet[i], plt.GetSubLegend(i))
			}
		}
	}

	if err := plot.Save(500, 500, config.Name); err != nil {
		return err
	}

	return nil
}

func initPlot(config *ExtraGraphicCfg) (*plot.Plot, error) {
	scaler := plot.LinearScale{}
	plot, err := plot.New()
	if err != nil {
		return nil, err
	}

	plot.X.Label.Text = config.XAsixLabel
	plot.Y.Label.Text = config.YAsixLabel
	plot.Title.Text = config.GraphicLabel

	if config.DoLegend {

		vertical := strings.Fields(config.LegendPosition)[0]
		horizontal := strings.Fields(config.LegendPosition)[1]
		if vertical == BOTTOMPOS {
			plot.Legend.Top = false
		} else {
			plot.Legend.Top = true
		}

		if horizontal == RIGHTPOS {
			plot.Legend.Left = false
		} else {
			plot.Legend.Left = true
		}
	}

	switch config.Type {
	case MULTYBARS:
		plot.X.Scale = scaler
	}

	return plot, nil
}

func initPlotter(config *ExtraGraphicCfg) (plot.Plotter, error) {
	switch config.Type {
	case LINEPLOT:
		pts, err := getXYSPoints(config)
		if err != nil {
			return nil, err
		}
		plt, err := plotter.NewLine(pts)
		if err != nil {
			return nil, err
		}
		plt.Color = red
		return plt, nil
	case CANDELPLOT:
		cnd, err := getCandels(config)
		if err != nil {
			return nil, err
		}

		plt, err := csplotter.NewCandlesticks(cnd)
		if err != nil {
			return nil, err
		}
		if len(config.CFG.ColorSet) >= 2 {
			plt.ColorUp, err = ParseHexColor(config.CFG.ColorSet[0])
			if err != nil {
				return nil, err
			}
			plt.ColorDown, err = ParseHexColor(config.CFG.ColorSet[1])
			if err != nil {
				return nil, err
			}
		} else {
			if strings.Contains(config.Flag, "inv") {
				plt.ColorUp = red
				plt.ColorDown = green
			} else {
				plt.ColorUp = green
				plt.ColorDown = red
			}
		}

		return plt, nil

	case MULTYBARS:
		bars, err := getMultyBars(config)
		if err != nil {
			return nil, err
		}

		colors := make([]color.Color, len(config.CFG.ColorSet))

		//get color set
		for i, cl := range config.CFG.ColorSet {
			colors[i], err = ParseHexColor(cl)
			if err != nil {
				return nil, err
			}
		}

		var length float64
		if strings.Contains(config.Flag, "length") {
			length, err = parseNumberFromFlag(config.Flag, "length")
			if err != nil {
				return nil, err
			}
		} else {
			length = 20
		}

		plt, err := mbplotter.NewMultiBarPlotter(bars, vg.Length(length), colors)
		if err != nil {
			return nil, err
		}

		return plt, nil
	default:
		return nil, UnknownOperation
	}
}

func createParamMap(names []string) map[string]int {
	m := make(map[string]int)
	for i, v := range names {
		m[v] = i
	}

	return m
}

func parseOperand(m map[string]int, str string, config *ExtraGraphicCfg, values []string) (float64, error) {
	if strings.Contains(str, "!") {
		opernum, err := strconv.Atoi(str[5:])
		if err != nil {
			return 0.0, err
		}
		oper := config.Operation[opernum]
		set := strings.Fields(oper)
		operand1, err := parseOperand(m, set[1], config, values)
		if err != nil {
			return 0.0, err
		}
		operand2, err := parseOperand(m, set[2], config, values)
		if err != nil {
			return 0.0, nil
		}
		switch set[0] {
		case SUB:
			return operand1 - operand2, nil
		case PLUS:
			return operand1 + operand2, nil
		case DIV:
			return operand1 / operand2, nil
		case TIME:
			return operand1 * operand2, nil
		default:
			return 0.0, UnknownOperation
		}
	} else {
		pos, ok := m[str]
		if !ok {
			return strconv.ParseFloat(str, 64)
		}
		op, err := strconv.ParseFloat(values[pos], 64)
		if err != nil {
			sub := strings.Trim(values[m[str]], "ms")
			op, err = strconv.ParseFloat(sub, 64)
			if err != nil {
				return 0.0, err
			}
		}
		return op, nil
	}
}

func parseNumberFromFlag(Flag, flg string) (float64, error) {
	var ymin float64
	index := strings.Index(Flag, flg+"=")
	if index == -1 {
		return 0.0, WrongFlagMB
	}
	index += len(flg) + 1
	substr := Flag[index:]
	index = strings.Index(substr, " ")
	if index != -1 {
		substr = substr[:index]
	}
	ymin, err := strconv.ParseFloat(substr, 64)
	if err != nil {
		return 0.0, WrongFlagMB
	}
	return ymin, nil
}
