package main

import (
	"math"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

type Calculator interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculatorService struct {
	points [][]float64
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{
		points: make([][]float64, 0),
	}
}

func (c *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(c.points) > 0 {
		prevPoints := c.points[len(c.points)-1]
		distance = calculateDistance(prevPoints[0], prevPoints[1], data.Lat, data.Long)
	}
	c.points = append(c.points, []float64{data.Lat, data.Long})
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
