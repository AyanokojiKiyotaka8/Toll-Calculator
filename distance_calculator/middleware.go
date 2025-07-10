package main

import (
	"time"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) CalculateDistance(data *types.OBUData) (distance float64) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"dist": distance,
			"took": time.Since(start),
		}).Info("calculate distance")
	}(time.Now())
	distance = l.next.CalculateDistance(data)
	return
}
