package service

import (
	"context"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/launchpals/open-now/backend/maps"
	open_now "github.com/launchpals/open-now/proto/go"
)

// Server is implements the open_now gRPC server
type Server struct {
	l    *zap.SugaredLogger
	core open_now.CoreServer

	m *maps.Client
}

// New instantiates a new server
func New(l *zap.SugaredLogger, m *maps.Client) (*Server, error) {
	return &Server{
		l: l,
		m: m,
	}, nil
}

// Run spins up the open_now service
func (s *Server) Run(ctx context.Context, host, port string) error {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on '%s:%s'", host, port)
	}

	// set logger to record all incoming requests
	grpcLogger := s.l.Desugar().Named("grpc")
	grpc_zap.ReplaceGrpcLogger(grpcLogger)
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Duration("grpc.duration", duration)
		}),
	}
	serverOpts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(grpcLogger, zapOpts...)),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(grpcLogger, zapOpts...)),
	}

	// initialize server
	gs := grpc.NewServer(serverOpts...)
	open_now.RegisterCoreServer(gs, s)

	// interrupt server gracefully if context is cancelled
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.l.Info("shutting down server")
				gs.GracefulStop()
				return
			}
		}
	}()

	// spin up server
	s.l.Infow("spinning up server",
		"host", host,
		"port", port)
	return gs.Serve(listener)
}

// GetStatus blah blah
func (s *Server) GetStatus(context.Context, *open_now.Empty) (*open_now.Status, error) {
	return &open_now.Status{}, nil
}

// GetPointsOfInterest blah blah
func (s *Server) GetPointsOfInterest(ctx context.Context, pos *open_now.Position) (*open_now.PointsOfInterest, error) {
	pois, err := s.m.PointsOfInterest(ctx, pos.GetCoordinates(), pos.GetSituation().GetSituation())
	if err != nil {
		// TODO: error handling
	}

	return &open_now.PointsOfInterest{
		Interests: pois,
	}, nil
}

// GetDirections blah blah
func (s *Server) GetDirections(context.Context, *open_now.DirectionsReq) (*open_now.DirectionsResp, error) {
	return nil, nil
}
