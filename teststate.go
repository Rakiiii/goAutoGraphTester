package main

import (
	"log"
	"time"
)

type TestState struct {
	itterator int
	maxItter  int
	maxTime   int64
	edgesMax  int
	vertexMax int
	testcond  string
}

func NewTestState(config *TestConfig) *TestState {
	return &TestState{itterator: -1,
		maxItter:  config.AmountOfItterations,
		maxTime:   config.MaxTimeForItteration,
		edgesMax:  -1,
		vertexMax: -1,
		testcond:  config.TypeOfStopCondition,
	}
}

func (t *TestState) Itterator() int {
	return t.itterator
}

func (t *TestState) isContinue(tm time.Duration, edg int, vertex int) bool {
	switch t.testcond {
	case TIMESTOP:
		return tm.Milliseconds() < t.maxTime
	case ITSTOP:
		return t.itterator < t.maxItter
	case MIXEDSTOP:
		return (tm.Milliseconds() < t.maxTime || t.itterator < t.maxItter)
	case EDGSTOP:
		return edg < t.edgesMax
	case VERTEXSTOP:
		return vertex < t.vertexMax
	default:
		log.Panic("Wrong stop condition:", t.testcond)
		return false
	}
}
