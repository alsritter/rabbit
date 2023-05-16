package server

import (
	"alsritter.icu/rabbit-template/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	// gRPC 服务暂时不使用
	srv := grpc.NewServer(opts...)
	// v1.RegisterDomainServiceServer(srv, domainService)
	// v1.RegisterUserServiceServer(srv, userService)
	// v1.RegisterManagerServiceServer(srv, managerService)
	// v1.RegisterOrderServiceServer(srv, orderService)
	// v1.RegisterTestServer(srv, testService)
	return srv
}
