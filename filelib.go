package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	//"math"
	"image/color"
	"errors"
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"

	csplotter "github.com/pplcc/plotext/custplotter"
	mbplotter "github.com/Rakiiii/goMultiBarPlotter"

)

const (
	PLOTCUSTNAME = "graphic.png"
	MARKDIFFNAME = "graphicDiff.png"
	MARKPROGRESSIONNAME = "graphicProgression.png"
	VERTEXTESTASICLABEL = "Amount of vertex"
	EDGETESTASICLABEL = "Amount of edges"
	ITTESTASICLABEL = "Amount of itterations"
	TIMEASICLABEL = "Time,ms"
	ADVTIMENAME = "AdvTimeGraphic.png"
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
	if conf.TypeOfTest != ITTEST {
		getPoints = func() (plotter.XYs, error) {
			pts := make(plotter.XYs, n)
			itter := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				elems := strings.Fields(scanner.Text())

				if conf.TypeOfTest == VERTEXTEST {
					num, err := strconv.Atoi(elems[0])
					if err != nil {
						return nil, err
					}
					pts[itter].X = float64(num)
				}

				if conf.TypeOfTest == EDGETEST {
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
	case VERTEXTEST:
		plot.X.Label.Text = VERTEXTESTASICLABEL
		plot.Y.Label.Text = TIMEASICLABEL
	case EDGETEST:
		plot.X.Label.Text = EDGETESTASICLABEL
		plot.Y.Label.Text = TIMEASICLABEL
	case ITTEST:
		plot.X.Label.Text = ITTESTASICLABEL
		plot.Y.Label.Text = TIMEASICLABEL

	}

	err = plotutil.AddLinePoints(plot, graphicalPoints)

	err = plot.Save(500, 500, PLOTCUSTNAME)
	if err != nil {
		return err
	}

	return nil
}


func DrawMarkDiff(file *os.File, conf *TestConfig, n int)error{

	var markpos int
	var respos int
	if conf.TypeOfTest != ITTEST {
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

				if conf.TypeOfTest == VERTEXTEST {
					ver, err := strconv.Atoi(elems[0])
					if err != nil {
						return nil, err
					}
					pts[itter].T = float64(ver)
				}

				if conf.TypeOfTest == EDGETEST {
					edg, err := strconv.Atoi(elems[1])
					if err != nil {
						return nil, err
					}
					pts[itter].T = float64(edg)
				}

				if conf.TypeOfTest != ITTEST {
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
				if num != 0{
					pts[itter].C = float64(num)
					pts[itter].H = float64(num)
				}else{
					pts[itter].C = float64(res)
					pts[itter].H = float64(res)
				}
				pts[itter].L = float64(res)
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
	case VERTEXTEST:
		plot.X.Label.Text = VERTEXTESTASICLABEL
	case EDGETEST:
		plot.X.Label.Text = EDGETESTASICLABEL
	case ITTEST:
		plot.X.Label.Text = ITTESTASICLABEL

	}

	plot.Y.Label.Text = "mark and result"

	bars, err := csplotter.NewCandlesticks(graphicalPoints)
	if err != nil {
		return err
	}

	bars.ColorUp = red
	bars.ColorDown = green 

	plot.Add(bars)

	err = plot.Save(500, 500, MARKDIFFNAME)
	if err != nil {
		return err
	}

	return nil
}

func DrawMarkProgression(file *os.File, conf *TestConfig, n int)error{

	var markpos int
	var respos int
	if conf.TypeOfTest != ITTEST {
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

				if conf.TypeOfTest == VERTEXTEST {
					ver, err := strconv.Atoi(elems[0])
					if err != nil {
						return nil, err
					}
					pts[itter].X = float64(ver)
				}

				if conf.TypeOfTest == EDGETEST {
					edg, err := strconv.Atoi(elems[1])
					if err != nil {
						return nil, err
					}
					pts[itter].X = float64(edg)
				}

				if conf.TypeOfTest != ITTEST {
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
				sub := num - res
				if sub >= 0{
					pts[itter].Y = float64(sub)
				}else{
					pts[itter].Y = 0.0
				}
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
	case VERTEXTEST:
		plot.X.Label.Text = VERTEXTESTASICLABEL
	case EDGETEST:
		plot.X.Label.Text = EDGETESTASICLABEL
	case ITTEST:
		plot.X.Label.Text = ITTESTASICLABEL

	}

	plot.Y.Label.Text = "Diff mark and result"


	err = plotutil.AddLinePoints(plot, graphicalPoints)

	
	err = plot.Save(500, 500, MARKPROGRESSIONNAME)
	if err != nil {
		return err
	}

	return nil
}

func DrawAdvtimeDistribution(file *os.File, config *TestConfig,n int)error{
	if config.ATCFG.EnableAdvTime{
		if len(config.ATCFG.GraphicCFG.NameSet) != len(config.ATCFG.GraphicCFG.ColorSet){
			return errors.New("Different length of color and name sets")
		}
		getPoints := func()(mbplotter.Bars,error){
			bars := make(mbplotter.Bars,n)
			scanner := bufio.NewScanner(file)
			it := 0
			for scanner.Scan(){
				bars[it].Ymin = 0.0

				elems := strings.Fields(scanner.Text())

				x,err := strconv.Atoi(elems[0])
				if err != nil{
					return nil , err
				}
				bars[it].X = float64(x)

				bars[it].Y = make([]float64,len(elems)-1)
				for i := range bars[it].Y{
					y,err := strconv.Atoi(elems[i+1])
					if err != nil{
						return nil , err
					}
					bars[it].Y[i] = float64(y)
				}
				
				//fmt.Println("X:",bars[it].X," Y:",bars[it].Y)	
				it++
			}
			return bars,nil
		}

		linearScale := plot.LinearScale{}

		bars,err := getPoints()
		if err != nil{
			return err
		}

		colors := make([]color.Color,len(config.ATCFG.GraphicCFG.ColorSet))
		//get color set
		for i,cl := range config.ATCFG.GraphicCFG.ColorSet{
			colors[i],err = ParseHexColor(cl)
			if err != nil{
				return err
			}
		}

		plotter,err := mbplotter.NewMultiBarPlotter(bars,vg.Length(10),colors)
		if err != nil{
			return err
		}

		//create new plot and set it settings
		plot, err := plot.New()
		if err != nil {
			return err
		}

		plot.Add(plotter)

		plot.Title.Text = "Advanced time distribution"
		plot.X.Scale = linearScale

		switch config.TypeOfTest {
		case VERTEXTEST:
			plot.X.Label.Text = VERTEXTESTASICLABEL
			plot.Y.Label.Text = TIMEASICLABEL
		case EDGETEST:
			plot.X.Label.Text = EDGETESTASICLABEL
			plot.Y.Label.Text = TIMEASICLABEL
		case ITTEST:
			plot.X.Label.Text = ITTESTASICLABEL
			plot.Y.Label.Text = TIMEASICLABEL
	
		}

		plot.Legend.Top = true
		plot.Legend.Left = true

		for i := range plotter.Colors{
			plot.Legend.Add(config.ATCFG.GraphicCFG.NameSet[i],plotter.GetSubLegend(i))
		}

		err = plot.Save(600, 600,ADVTIMENAME )
		if err != nil {
			return err
		}

		return nil
	}else{
		return nil
	}
}

func ParseHexColor(s string) (c color.RGBA, err error) {
    c.A = 0xff
    switch len(s) {
    case 7:
        _, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
    case 4:
        _, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
        // Double the hex digits:
        c.R *= 17
        c.G *= 17
        c.B *= 17
    default:
        err = fmt.Errorf("invalid length, must be 7 or 4")

    }
    return
}