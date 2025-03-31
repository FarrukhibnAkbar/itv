package tracer

import (
	"fmt"
	"itv/monorepo/movie_service/configs"
	"itv/monorepo/library/log"
	"itv/monorepo/library/tracer"

	"github.com/opentracing/opentracing-go"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
)

// Load ...
func Load(config *configs.Configuration, logger log.Factory) opentracing.Tracer {

	fmt.Println("tracer")
	metricsFactory := jexpvar.NewFactory(10) // 10 buckets for histograms
	tracer := tracer.Init(config.ServiceName, metricsFactory, logger)

	return tracer
}
