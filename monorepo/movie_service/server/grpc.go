package server

import (
	"fmt"
	"itv/monorepo/library/helper"
	"itv/monorepo/movie_service/configs"
	"itv/monorepo/movie_service/service"
	"log"
	"net"
	"runtime/debug"
	"time"

	pb "itv/monorepo/proto/movie_service"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Start will start grpc server
func Start(config *configs.Configuration, tracer opentracing.Tracer, configServer *service.MovieService) {
	myRecoveryHandler :=
		func(p any) (err error) {
			log.Println(string(debug.Stack()))
			log.Printf("Recovered from panic:: %v", p)
			return status.Errorf(codes.Internal, "panic triggered: %v", p)
		}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				otgrpc.OpenTracingStreamServerInterceptor(tracer),
				recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(myRecoveryHandler)),
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(myRecoveryHandler)),
				otgrpc.OpenTracingServerInterceptor(tracer),
			),
		),
	)
	pb.RegisterMovieServiceServer(grpcServer, configServer)

	cfg := configs.Config()
	fmt.Println("Starting grpc server on port:", config.RPCPort)
	fmt.Println("environment:", cfg.Environment)
	//listenting tcp rpcport
	lis, err := net.Listen("tcp", config.RPCPort)
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in Config file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		fmt.Println("listening tcp error: ", err)
	}

	fmt.Println("crm server running on port : ", config.RPCPort)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(fmt.Errorf("failed to serve: %w", err))
		}
	}()
}
