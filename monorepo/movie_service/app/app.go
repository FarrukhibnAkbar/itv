package app

import (
	"context"
	"fmt"

	"itv/monorepo/movie_service/configs"
	"itv/monorepo/movie_service/pkg/db"
	"itv/monorepo/movie_service/pkg/tracer"
	"itv/monorepo/movie_service/server"
	"itv/monorepo/movie_service/service"
	"itv/monorepo/movie_service/storage"
	"itv/monorepo/library/log"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loggerInit(config *configs.Configuration) log.Factory {
	fmt.Println("logInit")
	loggerForTracer, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)

	zapLogger := loggerForTracer.With(zap.String("service", config.ServiceName))
	logger := log.NewFactory(zapLogger)
	return logger
}

type App struct {
	engine *fx.App
}

// Start starts app with context spesified
func (a *App) Start(ctx context.Context) {
	a.engine.Start(ctx)
}

// Run starts the application, blocks on the signals channel, and then gracefully shuts the application down
func (a *App) Run() {
	a.engine.Run()
}

// New returns fx app
func New() App {

	engine := fx.New(
		fx.Provide(
			configs.Config,
			db.Init,
			loggerInit,
			tracer.Load,
			storage.NewMovie,
			service.New,
		),

		fx.Invoke(
			server.Start,
		),

		fx.WithLogger(
			func() fxevent.Logger {
				return fxevent.NopLogger
			},
		),
	)

	return App{engine: engine}
}
