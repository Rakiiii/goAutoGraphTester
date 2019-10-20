package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
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
	getPoints := func() (plotter.XYs, error) {
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

	plot.Y.Label.Text = "Time,ms"

	if conf.TypeOfTest == "vertextest" {
		plot.X.Label.Text = "Amount of vertex"
	}

	if conf.TypeOfTest == "edgestest" {
		plot.X.Label.Text = "Amount of edges"
	}

	err = plotutil.AddLinePoints(plot, graphicalPoints)

	err = plot.Save(500, 500, "graphic.png")
	if err != nil {
		return nil
	}

	return nil
}
