package service

import (
	"context"
	"fmt"
	libsLog "itv/monorepo/library/log"
	"itv/monorepo/movie_service/mappers"
	"itv/monorepo/movie_service/storage"
	pb "itv/monorepo/proto/movie_service"

	"github.com/opentracing/opentracing-go"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler ...
type MovieService struct {
	logger    libsLog.Factory
	tracer    opentracing.Tracer
	movieRepo storage.IMovieRepo
	pb.UnimplementedMovieServiceServer
	mapper mappers.AllMappers
}

// New ...
func New(logger libsLog.Factory, tracer opentracing.Tracer, movieRepo storage.IMovieRepo) *MovieService {
	fmt.Println("rpc handler new")
	movieHandler := &MovieService{
		logger:    logger,
		tracer:    tracer,
		movieRepo: movieRepo,
	}
	return movieHandler
}

// CreateMovie ...
func (b *MovieService) CreateMovie(ctx context.Context, req *pb.Movie) (*pb.Empty, error) {
	b.logger.For(ctx).Info("CreateMovie start")
	err := b.movieRepo.CreateMovie(ctx, b.mapper.ToMovieStruct(req))
	if err != nil {
		b.logger.For(ctx).Error("error in CreateMovie: ", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}
	b.logger.For(ctx).Info("CreateMovie finish")
	return &pb.Empty{}, nil
}

// GetMovies ...
func (b *MovieService) GetMovies(ctx context.Context, req *pb.Pagination) (*pb.MovieList, error) {
	b.logger.For(ctx).Info("GetMovies start")
	movies, err := b.movieRepo.GetMovies(ctx, int(req.Limit), int(req.Page))
	if err != nil {
		b.logger.For(ctx).Error("error in GetMovies: ", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}
	b.logger.For(ctx).Info("GetMovies finish")
	return b.mapper.ToMovieList(movies), nil
}

// GetMovieByID ...
func (b *MovieService) GetMovie(ctx context.Context, req *pb.MovieId) (*pb.Movie, error) {
	b.logger.For(ctx).Info("GetMovieByID start")
	movie, err := b.movieRepo.GetMovieByID(ctx, req.Id)
	if err != nil {
		b.logger.For(ctx).Error("error in GetMovieByID: ", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}
	b.logger.For(ctx).Info("GetMovieByID finish")
	return b.mapper.ToProtoMovie(movie), nil
}

// UpdateMovie ...
func (b *MovieService) UpdateMovie(ctx context.Context, req *pb.Movie) (*pb.Empty, error) {
	b.logger.For(ctx).Info("UpdateMovie start")
	err := b.movieRepo.UpdateMovie(ctx, b.mapper.ToMovieStruct(req))
	if err != nil {
		b.logger.For(ctx).Error("error in UpdateMovie: ", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}
	b.logger.For(ctx).Info("UpdateMovie finish")
	return &pb.Empty{}, nil
}

// DeleteMovie ...
func (b *MovieService) DeleteMovie(ctx context.Context, req *pb.MovieId) (*pb.Empty, error) {
	b.logger.For(ctx).Info("DeleteMovie start")
	err := b.movieRepo.DeleteMovie(ctx, req.Id)
	if err != nil {
		b.logger.For(ctx).Error("error in DeleteMovie: ", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}
	b.logger.For(ctx).Info("DeleteMovie finish")
	return &pb.Empty{}, nil
}