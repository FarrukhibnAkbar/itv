package storage

import (
	"context"
	"itv/monorepo/movie_service/entity"
	"itv/monorepo/movie_service/storage/postgres"

	"gorm.io/gorm"
)

type IMovieRepo interface {
	CreateMovie(ctx context.Context, movie *entity.Movie) error
	GetMovies(ctx context.Context, limit, page int) (*entity.MovieList, error)
	GetMovieByID(ctx context.Context, id string) (*entity.Movie, error)
	UpdateMovie(ctx context.Context, movie *entity.Movie) error
	DeleteMovie(ctx context.Context, id string) error
}

func NewMovie(db *gorm.DB) IMovieRepo {
	return postgres.NewMovie(db)
}
