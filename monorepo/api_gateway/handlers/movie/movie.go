package movie

import (
	"errors"
	"fmt"
	"itv/monorepo/api_gateway/dependencies"
	"itv/monorepo/api_gateway/entity"
	"itv/monorepo/api_gateway/mappers"
	"itv/monorepo/library/helper"
	"itv/monorepo/library/log"
	"itv/monorepo/library/utils"
	"itv/monorepo/proto/movie_service"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MovieHandler interface {
	CreateMovie(w http.ResponseWriter, r *http.Request)
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetMovieByID(w http.ResponseWriter, r *http.Request)
	DeleteMovie(w http.ResponseWriter, r *http.Request)
	UpdateMovie(w http.ResponseWriter, r *http.Request)
}

type movieHandler struct {
	logger  log.Factory
	mappers mappers.AllMappers
}

func New(logger log.Factory) MovieHandler {
	return &movieHandler{
		logger:  logger,
		mappers: mappers.AllMappers{},
	}
}

// CreateMovie godoc
// @Summary      Create movie
// @Description  Create movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        movie  body      entity.Movie  true  "Movie"
// @Success      200  {object}  utils.Response{data=utils.ResponseForExec}  "OK"  example:{"data": {"info": "Movie created successfully", "id": "123e4567-e89b-12d3-a456-426614174000"}}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /v1/new/movie [post]
// @Security     ApiKeyAuth
// @Security     BearerAuth
func (b *movieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	b.logger.For(r.Context()).Info("CreateMovie start")
	var body entity.Movie

	// Parse request body
	err := utils.BodyParser(r, &body)
	if err != nil {
		bodyErr := utils.HandleBodyParseError(err)
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> BodyParser",
			Source:  "api_gateway",
			ErrText: fmt.Sprintf("field: %v  message: %v", bodyErr.Field, bodyErr.Message),
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, fmt.Sprintf("field: %v  message: %v", bodyErr.Field, bodyErr.Message))
		return
	}

	body.Id = uuid.NewString()

	_, err = dependencies.MovieServiceClient().CreateMovie(r.Context(), b.mappers.ToMovieProto(&body))
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> CreateMovie",
			Source:  "api_gateway",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, "error in movieCliet.CreateMovie")
		return
	}

	utils.WriteJSONWithSuccess(w, r, utils.ResponseForExec{
		Info: "Movie created successfully",
		ID:   body.Id,
	})
	b.logger.For(r.Context()).Info("CreateCategoryHandler finished")
}

// GetMovies godoc
// @Summary      Get movies
// @Description  Get movies
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        page  query      int false  "Page"
// @Param        limit  query      int false "Limit"
// @Success      200  {object}  utils.Response{data=entity.MovieList}  "OK"  example:{"data": {"movies": [{"id": "123e4567-e89b-12d3-a456-426614174000", "title": "Inception", "description": "A mind-bending thriller", "year": 2010, "plot": "A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a CEO.", "director": "Christopher Nolan"}], "count": 1}}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /v1/movies [get]
func (b *movieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	b.logger.For(r.Context()).Info("GetMovies start")

	// Parse request pagenation
	page, limit, err := utils.ParsePagination(r)
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> Pagenation",
			Source:  "api_gateway",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, err.Error())
		return
	}

	// Call grpc service
	movies, err := dependencies.MovieServiceClient().GetMovies(r.Context(), &movie_service.Pagination{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> GetMovies",
			Source:  "api_gateway",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, "error in movieCliet.GetMovies")
		return
	}
	if movies == nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> GetMovies",
			Source:  "api_gateway",
			ErrText: "no movies found",
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, "no movies found")
		return
	}

	utils.WriteJSONWithSuccess(w, r, b.mappers.ToMovieList(movies))

	b.logger.For(r.Context()).Info("GetMovies finished")

}

// GetMovieByID godoc
// @Summary      Get movie by id
// @Description  Get movie by id
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Movie ID"
// @Success      200  {object}  utils.Response{data=entity.Movie}  "OK"  example:{"data": {"id": "123e4567-e89b-12d3-a456-426614174000", "title": "Inception", "description": "A mind-bending thriller", "year": 2010, "plot": "A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a CEO.", "director": "Christopher Nolan"}}sk of planting an idea into the mind of a CEO.", "director": "Christopher Nolan"}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /v1/movie/{id} [get]
func (b *movieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	b.logger.For(r.Context()).Info("GetMovieByID start")

	vars := mux.Vars(r)
	id := vars["id"]

	if ok := utils.IsValidUUID(id); !ok {
		utils.HandleGrpcErrWithMessage(w, r, errors.New("invalid id"), "invalid id")
		return
	}

	// Call grpc service
	movie, err := dependencies.MovieServiceClient().GetMovie(r.Context(), &movie_service.MovieId{
		Id: id,
	})
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> GetMovieByID",
			Source:  "api_gateway",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, "error in movieCliet.GetMovieByID")
		return
	}
	if movie == nil {
		utils.HandleGrpcErrWithMessage(w, r, err, "no movies found")
		return
	}

	utils.WriteJSONWithSuccess(w, r, b.mappers.ToMovieStruct(movie))

	b.logger.For(r.Context()).Info("GetMovies finished")
}

// DeleteMovie godoc
// @Summary      Delete movie
// @Description  Delete movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Movie ID"  example(123e4567-e89b-12d3-a456-426614174000)
// @Success      200  {object}  utils.Response{data=utils.ResponseForExec}  "OK"  example:{"data": {"info": "Movie deleted successfully", "id": "123e4567-e89b-12d3-a456-426614174000"}}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /v1/movie/{id} [delete]
// @Security     ApiKeyAuth
// @Security     BearerAuth
func (b *movieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	b.logger.For(r.Context()).Info("DeleteMovie start")

	id := mux.Vars(r)["id"]
	if ok := utils.IsValidUUID(id); !ok {
		utils.HandleGrpcErrWithMessage(w, r, errors.New("invalid id"), "invalid id")
		return
	}

	// call grpc service
	_, err := dependencies.MovieServiceClient().DeleteMovie(r.Context(), &movie_service.MovieId{
		Id: id,
	})
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> DeleteMovie",
			Source:  "api_gateway",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, "error in movieCliet.DeleteMovie")
		return
	}
	utils.WriteJSONWithSuccess(w, r, utils.ResponseForExec{
		Info: "Movie deleted successfully",
		ID:   id,
	})
	b.logger.For(r.Context()).Info("DeleteMovie finished")
}

// UpdateMovie godoc
// @Summary      Update movie
// @Description  Update movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        movie  body      entity.Movie  true  "Movie"
// @Success      200  {object}  utils.Response{data=utils.ResponseForExec}  "OK"  example:{"data": {"info": "Movie updated successfully", "id": "123e4567-e89b-12d3-a456-426614174000"}}
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /v1/movie [put]
// @Security     ApiKeyAuth
// @Security     BearerAuth
func (b *movieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	b.logger.For(r.Context()).Info("UpdateMovie start")

	var body entity.Movie

	// Parse request body
	err := utils.BodyParser(r, &body)
	if err != nil {
		bodyErr := utils.HandleBodyParseError(err)
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> BodyParser",
			Source:  "api_gateway",
			ErrText: fmt.Sprintf("field: %v  message: %v", bodyErr.Field, bodyErr.Message),
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, fmt.Sprintf("field: %v  message: %v", bodyErr.Field, bodyErr.Message))
		return
	}

	_, err = dependencies.MovieServiceClient().UpdateMovie(r.Context(), b.mappers.ToMovieProto(&body))
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "movie-> UpdateMovie",
			Source:  "api_gateway",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		utils.HandleGrpcErrWithMessage(w, r, err, "error in movieCliet.UpdateMovie")
		return
	}

	utils.WriteJSONWithSuccess(w, r, utils.ResponseForExec{
		Info: "Movie updated successfully",
		ID:   body.Id,
	})
	b.logger.For(r.Context()).Info("UpdateCategoryHandler finished")
}
