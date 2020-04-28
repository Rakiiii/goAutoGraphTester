package main

import (
	"fmt"
	"image/color"
	"io"
	"os"
)

const (
	PLOTCUSTNAME        = "graphic.png"
	MARKDIFFNAME        = "graphicDiff.png"
	MARKPROGRESSIONNAME = "graphicProgression.png"
	VERTEXTESTASICLABEL = "Amount of vertex"
	EDGETESTASICLABEL   = "Amount of edges"
	ITTESTASICLABEL     = "Amount of itterations"
	TIMEASICLABEL       = "Time,ms"
	ADVTIMENAME         = "AdvTimeGraphic.png"
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

func addToAdvTime(line []string) error {
	if err := advtimewriter.Write(line); err != nil {
		return err
	}
	return nil
}

func addToResult(line []string) error {
	if err := resultwriter.Write(line); err != nil {
		return err
	}
	return nil
}

func flushWriters() {
	advtimewriter.Flush()
	resultwriter.Flush()
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
