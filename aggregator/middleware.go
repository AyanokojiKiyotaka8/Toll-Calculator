package main

import (
	"time"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) AggregateDistance(dist *types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(start),
			"error": err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = l.next.AggregateDistance(dist)
	return
}

func (l *LogMiddleware) CalculateInvoice(id int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		totalAmount := 0.0
		totalDistance := 0.0
		if invoice != nil {
			totalAmount = invoice.TotalAmount
			totalDistance = invoice.TotalDistance
		}
		logrus.WithFields(logrus.Fields{
			"took":          time.Since(start),
			"totalDistance": totalDistance,
			"totalAmount":   totalAmount,
			"obuID":         id,
			"error":         err,
		}).Info("CalculateInvoice")
	}(time.Now())
	invoice, err = l.next.CalculateInvoice(id)
	return
}
