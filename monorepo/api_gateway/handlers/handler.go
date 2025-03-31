package handlers

import (
	"itv/monorepo/api_gateway/handlers/movie"
	"itv/monorepo/library/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Handlers struct {
	MovieHandlers movie.MovieHandler
}

func New() (*Handlers, error) {
	zapLogger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)
	logger := log.NewFactory(zapLogger)
	return &Handlers{
		MovieHandlers: movie.New(logger),
	}, nil
}
