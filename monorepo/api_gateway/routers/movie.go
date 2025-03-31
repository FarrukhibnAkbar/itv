package routers

import (
	"itv/monorepo/api_gateway/handlers"
	"itv/monorepo/api_gateway/pkg/tracing"
	"net/http"
)

func RegisterMovieRoutes(r *tracing.TracedServeMux, h *handlers.Handlers) {
	// for movie
	r.Handle("POST", "/v1/movie", http.HandlerFunc(h.MovieHandlers.CreateMovie))
	r.Handle("GET", "/v1/movies", http.HandlerFunc(h.MovieHandlers.GetMovies))
	r.Handle("GET", "/v1/movie/{id}", http.HandlerFunc(h.MovieHandlers.GetMovieByID))
	r.Handle("DELETE", "/v1/movie/{id}", http.HandlerFunc(h.MovieHandlers.DeleteMovie))
	r.Handle("PUT", "/v1/movie", http.HandlerFunc(h.MovieHandlers.UpdateMovie))
}
