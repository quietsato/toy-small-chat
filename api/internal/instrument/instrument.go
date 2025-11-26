package instrument

import (
	"context"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func setupLogger(logLevel slog.Level) {
	customLogger := slog.New(
		slogmulti.Fanout(
			otelslog.NewHandler("toy-small-chat"),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
			}),
		),
	)
	slog.SetDefault(customLogger)
}

func Init(ctx context.Context, logLevel slog.Level, otlpEndpoint string) func() {
	setupLogger(logLevel)

	resource, _ := NewResource()
	shutdown, err := InitTracer(ctx, resource, otlpEndpoint)
	if err != nil {
		slog.Error("failed to initialize tracer", slog.Any("err", err))
	}

	loggerProvider, err := InitLogger(ctx, resource, otlpEndpoint)
	if err != nil {
		slog.Error("failed to initialize logger", slog.Any("err", err))
	}

	return func() {
		defer func() {
			if err := shutdown(ctx); err != nil {
				slog.WarnContext(ctx, "failed to shutdown tracer", slog.Any("err", err))
			}
		}()
		defer func() {
			if err := loggerProvider.Shutdown(ctx); err != nil {
				slog.WarnContext(ctx, "failed to shutdown logger provider", slog.Any("err", err))
			}
		}()
	}
}
