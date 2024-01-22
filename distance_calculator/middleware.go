package main

import (
	"time"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Calculator
}

func NewLogMiddleware(next Calculator) Calculator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.OBUData) (distance float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"distance": distance,
		}).Info("calculate distance")
	}(time.Now())
	distance, err = m.next.CalculateDistance(data)
	return
}
