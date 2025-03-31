package mappers

import (
	e "itv/monorepo/api_gateway/entity"
	b "itv/monorepo/proto/movie_service"
)

type AllMappers struct{}

func (AllMappers) ToMovieProto(movie *e.Movie) *b.Movie {
	return &b.Movie{
		Id:          movie.Id,
		Title:       movie.Title,
		Description: movie.Description,
		Director:    movie.Director,
		Year:        movie.Year,
		Plot:        movie.Plot,
	}
}

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

func (AllMappers) ToMovieList(movies *b.MovieList) *e.MovieList {
	movieList := make([]e.Movie, len(movies.Movies))
	for i, movie := range movies.Movies {
		movieList[i] = e.Movie{
			Id:          movie.Id,
			Title:       movie.Title,
			Description: movie.Description,
			Director:    movie.Director,
			Year:        movie.Year,
			Plot:        movie.Plot,
		}
	}
	return &e.MovieList{
		Movies: movieList,
		Count:  movies.Count,
	}
}
