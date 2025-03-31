package mappers

import (
	e "itv/monorepo/movie_service/entity"
	b "itv/monorepo/proto/movie_service"
)

type AllMappers struct{}

func (AllMappers) ToMovieStruct(movie *b.Movie) *e.Movie {
	return &e.Movie{
		Id:          movie.Id,
		Title:       movie.Title,
		Description: movie.Description,
		Director:    movie.Director,
		Year:        movie.Year,
		Plot:        movie.Plot,
	}
}

func (AllMappers) ToMovieList(movies *e.MovieList) *b.MovieList {
	movieList := make([]*b.Movie, len(movies.Movies))
	for i, movie := range movies.Movies {
		movieList[i] = &b.Movie{
			Id:          movie.Id,
			Title:       movie.Title,
			Description: movie.Description,
			Director:    movie.Director,
			Year:        movie.Year,
			Plot:        movie.Plot,
		}
	}
	return &b.MovieList{
		Movies: movieList,
		Count:  movies.Count,
	}
}

func (AllMappers) ToProtoMovie(movie *e.Movie) *b.Movie {
	return &b.Movie{
		Id:          movie.Id,
		Title:       movie.Title,
		Description: movie.Description,
		Director:    movie.Director,
		Year:        movie.Year,
		Plot:        movie.Plot,
	}
}