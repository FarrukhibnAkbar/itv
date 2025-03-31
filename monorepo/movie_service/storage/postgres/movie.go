package postgres

import (
	"context"
	"itv/monorepo/library/helper"
	"itv/monorepo/library/utils"
	"itv/monorepo/movie_service/constants"
	"itv/monorepo/movie_service/entity"
	"time"

	_ "github.com/lib/pq" //db driver
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type movieRepo struct {
	db *gorm.DB
}

// NewMovie ...
func NewMovie(db *gorm.DB) *movieRepo {
	return &movieRepo{db: db}
}

// CreateMovie - this upserts the movie function
func (r *movieRepo) CreateMovie(ctx context.Context, movie *entity.Movie) error {
	err := r.db.WithContext(ctx).Table("movie").Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "description", "director", "year", "plot"}),
		}).Create(movie).Error
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/storage/postgres/movie.go file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		return err
	}
	return nil
}

// GetMovies - this gets all the movies
func (r *movieRepo) GetMovies(ctx context.Context, limit, page int) (*entity.MovieList, error) {
	var movies []entity.Movie
	var count int64
	offset := limit * (page - 1)
	err := r.db.WithContext(ctx).Table(constants.MovieTableName).
		Where("state=?", constants.Active).Offset(offset).Limit(limit).Find(&movies).Error
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/storage/postgres/movie.go file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		return nil, err
	}

	err = r.db.WithContext(ctx).Table(constants.MovieTableName).Where("state=?", constants.Active).Count(&count).Error
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/storage/postgres/movie.go file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		return nil, err
	}
	if count == 0 {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/storage/postgres/movie.go file",
			Source:  "movie_service",
			ErrText: "no movies found",
			Time:    time.Now().Format(time.RFC3339),
		})
		return nil, utils.HandleDBError("GetMovies", err, r.db)
	}
	return &entity.MovieList{
		Movies: movies,
		Count:  int32(count),
	}, nil
}

// GetMovieByID - this gets a movie by id
func (r *movieRepo) GetMovieByID(ctx context.Context, id string) (*entity.Movie, error) {
	var movie entity.Movie
	err := r.db.WithContext(ctx).Table(constants.MovieTableName).Where("id=? and state=?", id, constants.Active).First(&movie).Error
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/storage/postgres/movie.go file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		return nil, err
	}
	return &movie, nil
}

// UpdateMovie - this updates a movie
func (r *movieRepo) UpdateMovie(ctx context.Context, movie *entity.Movie) error {
	err := r.db.WithContext(ctx).Table(constants.MovieTableName).Where("id=?", movie.Id).Updates(movie).Error
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/storage/postgres/movie.go file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		return err
	}
	return nil
}

// DeleteMovie - this deletes a movie
func (r *movieRepo) DeleteMovie(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Table(constants.MovieTableName).Where("id=?", id).Update("state", constants.InActive).Error
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/storage/postgres/movie.go file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		return utils.HandleDBError("DeleteMovie", err, r.db)
	}
	return nil
}
