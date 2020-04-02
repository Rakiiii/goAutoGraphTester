package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"

	csplotter "github.com/pplcc/plotext/custplotter"

)

var (
	red   = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green = color.RGBA{R: 0, G: 200, B: 100, A: 255}
)

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func AppendStringToFile(path, text string, it int) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	if it != 0 {
		_, err = f.WriteString("\n" + text)
	} else {
		_, err = f.WriteString(text)
	}
	if err != nil {
		return err
	}
	return nil
}

func DrawPlotCust(file *os.File, conf *TestConfig, n int) error {
	var getPoints func() (plotter.XYs, error)
	if conf.TypeOfTest != "ittest" {
		getPoints = func() (plotter.XYs, error) {
			pts := make(plotter.XYs, n)
			itter := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				elems := strings.Fields(scanner.Text())

				if conf.TypeOfTest == "vertextest" {
					num, err := strconv.Atoi(elems[0])
					if err != nil {
						return nil, err
					}
					pts[itter].X = float64(num)
				}

				if conf.TypeOfTest == "edgestest" {
					num, err := strconv.Atoi(elems[1])
					if err != nil {
						return nil, err
					}
					pts[itter].X = float64(num)
				}

				num, err := strconv.Atoi(strings.Trim(elems[4], "ms"))
				if err != nil {
					return nil, err
				}

				pts[itter].Y = float64(num)
				itter++
			}
			return pts, scanner.Err()
		}
	} else {
		getPoints = func() (plotter.XYs, error) {
			pts := make(plotter.XYs, n)
			itter := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				elems := strings.Fields(scanner.Text())
				num, err := strconv.Atoi(elems[0])
				if err != nil {
					return nil, err
				}
				pts[itter].X = float64(num)

				num, err = strconv.Atoi(elems[1])
				if err != nil {
					return nil, err
				}

				pts[itter].Y = float64(num)
				itter++
			}
			return pts, scanner.Err()
		}

	}
	graphicalPoints, err := getPoints()
	if err != nil {
		return err
	}

	//create new plot and set it settings
	plot, err := plot.New()
	if err != nil {
		return err
	}

	plot.Title.Text = conf.GraphicTitle

	switch conf.TypeOfTest {
	case "vertextest":
		plot.X.Label.Text = "Amount of vertex"
		plot.Y.Label.Text = "Time,ms"
	case "edgestest":
		plot.X.Label.Text = "Amount of edges"
		plot.Y.Label.Text = "Time,ms"
	case "ittest":
		plot.X.Label.Text = "Amount of itterations"
		plot.Y.Label.Text = "Value"

	}

	err = plotutil.AddLinePoints(plot, graphicalPoints)

	err = plot.Save(500, 500, "graphic.png")
	if err != nil {
		return nil
	}

	return nil
}


func DrawMarkDiff(file *os.File, conf *TestConfig, n int)error{

	var markpos int
	var respos int
	if conf.TypeOfTest != "ittest" {
		markpos = 6
		respos = 5
	}else{
		markpos = 5
		respos = 1 
	}
	getPoints := func() (csplotter.TOHLCVs, error) {
			pts := make(csplotter.TOHLCVs, n)
			itter := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				elems := strings.Fields(scanner.Text())

				if conf.TypeOfTest == "vertextest" {
					ver, err := strconv.Atoi(elems[0])
					if err != nil {
						return nil, err
					}
					pts[itter].T = float64(ver)
				}

				if conf.TypeOfTest == "edgestest" {
					edg, err := strconv.Atoi(elems[1])
					if err != nil {
						return nil, err
					}
					pts[itter].T = float64(edg)
				}

				if conf.TypeOfTest != "ittest" {
					it, err := strconv.Atoi(elems[0])
					if err != nil {
						return nil, err
					}
					pts[itter].T = float64(it)
				}


				res, err := strconv.Atoi(elems[respos])
				if err != nil {
					return nil, err
				}

				num, err := strconv.Atoi(elems[markpos])
				if err != nil {
					return nil, err
				}
				pts[itter].O = float64(res)
				pts[itter].C = float64(num)
				pts[itter].L = float64(res)
				pts[itter].H = float64(num)
				itter++
			}
			return pts, scanner.Err()
	}

	graphicalPoints, err := getPoints()
	if err != nil {
		return err
	}

	//create new plot and set it settings
	plot, err := plot.New()
	if err != nil {
		return err
	}


	plot.Title.Text = "mark and result"

	switch conf.TypeOfTest {
	case "vertextest":
		plot.X.Label.Text = "Amount of vertex"
	case "edgestest":
		plot.X.Label.Text = "Amount of edges"
	case "ittest":
		plot.X.Label.Text = "Amount of itterations"

	}

	plot.Y.Label.Text = "mark and result"

	bars, err := csplotter.NewCandlesticks(graphicalPoints)
	if err != nil {
		return err
	}

	bars.ColorUp = red
	bars.ColorDown = green 

	plot.Add(bars)

	err = plot.Save(500, 500, "graphicDiff.png")
	if err != nil {
		return nil
	}

	return nil
}

func DrawMarkProgression(file *os.File, conf *TestConfig, n int)error{

	var markpos int
	var respos int
	if conf.TypeOfTest != "ittest" {
		markpos = 6
		respos = 5
	}else{
		markpos = 5
		respos = 1 
	}
	getPoints := func() (plotter.XYs, error) {
			pts := make(plotter.XYs, n)
			itter := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				elems := strings.Fields(scanner.Text())

				if conf.TypeOfTest == "vertextest" {
					ver, err := strconv.Atoi(elems[0])
					if err != nil {
						return nil, err
					}
					pts[itter].X = float64(ver)
				}

				if conf.TypeOfTest == "edgestest" {
					edg, err := strconv.Atoi(elems[1])
					if err != nil {
						return nil, err
					}
					pts[itter].X = float64(edg)
				}

				if conf.TypeOfTest != "ittest" {
					it, err := strconv.Atoi(elems[0])
					if err != nil {
						return nil, err
					}
					pts[itter].X = float64(it)
				}


				res, err := strconv.Atoi(elems[respos])
				if err != nil {
					return nil, err
				}

				num, err := strconv.Atoi(elems[markpos])
				if err != nil {
					return nil, err
				}
				pts[itter].Y = float64(num - res)
				itter++

				
			}
			return pts, scanner.Err()
	}

	
	graphicalPoints, err := getPoints()
	if err != nil {
		return err
	}

	//create new plot and set it settings
	plot, err := plot.New()
	if err != nil {
		return err
	}


	plot.Title.Text = "progression of mark and result difference"

	switch conf.TypeOfTest {
	case "vertextest":
		plot.X.Label.Text = "Amount of vertex"
	case "edgestest":
		plot.X.Label.Text = "Amount of edges"
	case "ittest":
		plot.X.Label.Text = "Amount of itterations"

	}

	plot.Y.Label.Text = "Diff mark and result"


	err = plotutil.AddLinePoints(plot, graphicalPoints)

	
	err = plot.Save(500, 500, "graphicProgression.png")
	if err != nil {
		return nil
	}

	return nil
}
