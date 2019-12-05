package lib

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func CreateTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	var cfg config.Configuration
	jLogger := log.StdLogger
	jMetricsFactory := metrics.NullFactory
	cfg.ServiceName = serviceName
	return cfg.NewTracer(
		config.Logger(jLogger),
		config.Metrics(jMetricsFactory),
	)
}

func InitGlobalTracer(serviceName string) (io.Closer, error) {

	return initGlobalTracer(config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}, serviceName)
}

func InitGlobalTracerProduction(serviceName string) (io.Closer, error) {
	return initGlobalTracer(config.Configuration{}, serviceName)
}

func initGlobalTracer(cfg config.Configuration, serviceName string) (io.Closer, error) {
	jLogger := log.StdLogger
	jMetricsFactory := metrics.NullFactory
	return cfg.InitGlobalTracer(
		serviceName,
		config.Logger(jLogger),
		config.Metrics(jMetricsFactory),
	)
}
