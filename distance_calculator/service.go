package main

import (
	"fmt"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

type Calculator interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculatorService struct{}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (c *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	fmt.Println("Calculating distance")
	return 00, nil
}
