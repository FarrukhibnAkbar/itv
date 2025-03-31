package dependencies

import (
	"itv/monorepo/api_gateway/configs"
	helper "itv/monorepo/library/helper"
	"itv/monorepo/proto/movie_service"
	"fmt"
	"sync"
	"time"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	movieServiceClient movie_service.MovieServiceClient
	onceMovieService   sync.Once
)

// AuthServiceClient loads AuthServiceClient using atomic pattern
func MovieServiceClient() movie_service.MovieServiceClient {
	onceMovieService.Do(func() {
		movieServiceClient = loadMovieServiceClient()
	})
	return movieServiceClient
}

func loadMovieServiceClient() movie_service.MovieServiceClient {
	tracer := opentracing.GlobalTracer()
	conf := configs.Config()
	connMovie, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.MovieServiceHost, conf.MovieServicePort),
		grpc.WithTransportCredentials(
			insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(tracer),
		),
	)
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie service grpc connection dial",
			Source:  "api_gateway",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		panic(fmt.Errorf("movie service dial host: %s port:%d err: %s",
			conf.MovieServiceHost, conf.MovieServicePort, err))
	}

	return movie_service.NewMovieServiceClient(connMovie)
}
