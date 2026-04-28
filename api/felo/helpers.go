package felo

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func unaryMethod[Req any, Resp any](name string, handler func(context.Context, *Req) (*Resp, error)) grpc.MethodDesc {
	return grpc.MethodDesc{
		MethodName: name,
		Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			req := new(Req)
			if err := dec(req); err != nil {
				return nil, err
			}
			if interceptor == nil {
				return handler(ctx, req)
			}
			info := &grpc.UnaryServerInfo{
				Server:     srv,
				FullMethod: name,
			}
			return interceptor(ctx, req, info, func(ctx context.Context, req any) (any, error) {
				typedReq, ok := req.(*Req)
				if !ok {
					return nil, fmt.Errorf("invalid request type")
				}
				return handler(ctx, typedReq)
			})
		},
	}
}
