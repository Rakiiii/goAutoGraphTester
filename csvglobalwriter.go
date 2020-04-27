package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

var (
	advtimewriter *custWriter
	resultwriter  *custWriter
)

type custWriter struct {
	file   *os.File
	writer *csv.Writer
}

func NewCustWriterF(file *os.File) *custWriter {
	return &custWriter{file: file, writer: csv.NewWriter(file)}
}

func NewCustWriter(path string) *custWriter {
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	return &custWriter{file: file, writer: csv.NewWriter(file)}
}

func (c *custWriter) Write(line []string) error {
	fmt.Println(line)
	return c.writer.Write(line)
}

func (c *custWriter) Flush() {
	c.writer.Flush()
	c.file.Close()
}
