package main

import (
	"context"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		logrus.WithFields(logrus.Fields{
			"obuID": dist.OBUID,
			"value": dist.Value,
			"unix":  dist.Unix,
			"error": err,
		}).Error("Failed to aggregate distance")
		return nil, status.Errorf(codes.Internal, "aggregation failed")
	}

	return &types.AggregatorResp{}, nil
}
