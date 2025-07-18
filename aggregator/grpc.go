package main

import (
	"context"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
)

type GRPCAggregatorServer struct {
	svc Aggregator
	types.UnimplementedAggregatorServer
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregatorReq) (*types.AggregatorResp, error) {
	dist := &types.Distance{
		OBUID: int(req.GetObuId()),
		Value: req.GetValue(),
		Unix:  req.GetUnix(),
	}
	if err := s.svc.AggregateDistance(dist); err != nil {
		return nil, err
	}
	return &types.AggregatorResp{}, nil
}
