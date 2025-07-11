package main

import "github.com/AyanokojiKiyotaka8/Toll-Calculator/types"

const roadPrice = 3.15

type Aggregator interface {
	AggregateDistance(*types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(dist *types.Distance) error {
	return i.store.Insert(dist)
}

func (i *InvoiceAggregator) CalculateInvoice(id int) (*types.Invoice, error) {
	dist, err := i.store.Get(id)
	if err != nil {
		return nil, err
	}
	return &types.Invoice{
		OBUID:         id,
		TotalDistance: dist,
		TotalAmount:   dist * roadPrice,
	}, nil
}
