package service

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	open_now "github.com/lunchpals/open-now/proto/go"
)

type Server struct {
	l *zap.SugaredLogger
	core open_now.CoreServer
}

func NewServer(l *zap.SugaredLogger) (*Server, error) {
	return &Server{l: l}, nil
}

func (s *Server) Run(host, port string) error {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return err
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
	pb.RegisterCoreServer(gs, s)

	// interrupt server gracefully if context is cancelled
	go func() {
		for {
			select {
			case <-ctx.Done():
				d.l.Info("shutting down server")
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
