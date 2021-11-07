package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jConf "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

type zapWrapper struct {
	logger *zap.Logger
}

// Error logs a message at error priority
func (w *zapWrapper) Error(msg string) {
	w.logger.Error(msg)
}

// Infof logs a message at info priority
func (w *zapWrapper) Infof(msg string, args ...interface{}) {
	w.logger.Sugar().Infof(msg, args...)
}

func InitJaeger(agentHost, service string, logger *zap.Logger) (opentracing.Tracer, io.Closer, error) {
	cfg := &jConf.Configuration{
		ServiceName: service,
		Sampler: &jConf.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jConf.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentHost,
		},
	}

	tracer, closer, err := cfg.NewTracer(jConf.Logger(&zapWrapper{logger: logger}))
	if err != nil {
		return nil, nil, err
	}

	return tracer, closer, nil
}
