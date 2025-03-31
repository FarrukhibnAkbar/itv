package main

import(
	"context"
	"fmt"
	"itv/monorepo/api_gateway/configs"
	"itv/monorepo/api_gateway/handlers"
	"itv/monorepo/api_gateway/middleware"
	"itv/monorepo/api_gateway/pkg/tracing"
	"itv/monorepo/api_gateway/routers"
	"itv/monorepo/library/tracer"
	"itv/monorepo/library/log"
	"itv/monorepo/library/utils"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"

	_ "time/tzdata"

	// sw "github.com/RussellLuo/slidingwindow"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "itv/monorepo/api_gateway/docs" // Import Swagger docs
    httpSwagger "github.com/swaggo/http-swagger"
)

// @title Movie API
// @version 1.0
// @description This is the API documentation for the Movie Service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func NewMux(lc fx.Lifecycle, conf *configs.Configuration) *tracing.TracedServeMux {
	metricsFactory := jexpvar.NewFactory(10) // 10 buckets for histograms
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)

	zapLogger := logger.With(zap.String("service", "api_gateway"))
	tracer := tracer.Init("api_gateway", metricsFactory, log.NewFactory(zapLogger))
	opentracing.SetGlobalTracer(tracer)

	root := mux.NewRouter()
	root = root.PathPrefix("/api").Subrouter()

	root.Use(middleware.PanicRecovery)
	root.Use(middleware.Logging)

	// Add Swagger route
    root.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	handlerWithTracer := tracing.NewServeMux(tracer, root)

	server := &http.Server{
		Addr:    conf.HTTPPort,
		Handler: middleware.Cors(handlerWithTracer),
		// ReadTimeout: , //TODO: should be set
		// WriteTimeout: //TODO: should be set
	}
	lc.Append(fx.Hook{

		OnStart: func(context.Context) error {
			fmt.Println("Starting HTTP server", server.Addr)
			fmt.Println("environment:", conf.Environment)
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return handlerWithTracer
}

func main() {

	utils.PrintMemStats()

	app := fx.New(
		fx.Provide(
			configs.Config,
			NewMux,
			handlers.New,
		),

		fx.Invoke(
			routers.RegisterMovieRoutes,
		),

	)

	utils.PrintMemStats()
	app.Run()
}