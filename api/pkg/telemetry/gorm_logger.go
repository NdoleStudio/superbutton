package telemetry

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

type gormLogger struct {
	tracer Tracer
	logger Logger
}

// NewGormLogger creates a new instance of gormLogger
func NewGormLogger(tracer Tracer, logger Logger) logger.Interface {
	return &gormLogger{
		tracer: tracer,
		logger: logger.WithService(fmt.Sprintf("%T", gormLogger{})),
	}
}

// LogMode log mode
func (gorm *gormLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return gorm
}

func (gorm *gormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	gorm.logger.WithSpan(gorm.tracer.Span(ctx).SpanContext()).Info(fmt.Sprintf(s, i...))
}

func (gorm *gormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	gorm.logger.WithSpan(gorm.tracer.Span(ctx).SpanContext()).Warn(fmt.Errorf(s, i...))
}

func (gorm *gormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	gorm.logger.WithSpan(gorm.tracer.Span(ctx).SpanContext()).Error(fmt.Errorf(s, i...))
}

func (gorm *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	l := gorm.logger.WithSpan(gorm.tracer.Span(ctx).SpanContext()).WithString("latency", elapsed.String())
	sql, rows := fc()
	l.Debug(fmt.Sprintf("[ROWS:%d][%s]", rows, sql))
}
