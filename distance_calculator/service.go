package main

import (
	"math"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
)

type CalculatorServicer interface {
	CalculateDistance(*types.OBUData) float64
}

type CalculatorService struct {
	prevLat     float64
	prevLong    float64
	initialized bool
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (c *CalculatorService) CalculateDistance(data *types.OBUData) float64 {
	if !c.initialized {
		c.prevLat = data.Lat
		c.prevLong = data.Long
		c.initialized = true
	}
	dist := c.computeDistance(c.prevLat, c.prevLong, data.Lat, data.Long)
	c.prevLat = data.Lat
	c.prevLong = data.Long
	return dist
}

func (c *CalculatorService) computeDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
