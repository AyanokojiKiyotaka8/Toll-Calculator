package client

import (
	"context"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregatorReq) error
}
