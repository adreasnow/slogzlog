package slogzlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func setupLogger(t *testing.T) (context.Context, *bytes.Buffer) {
	t.Helper()

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	buf := new(bytes.Buffer)
	return log.
		Output(
			io.MultiWriter(
				zerolog.TestWriter{T: t},
				zerolog.ConsoleWriter{Out: os.Stdout},
				zerolog.ConsoleWriter{Out: buf, NoColor: true},
			),
		).
		WithContext(t.Context()), buf
}

func TestHandler(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		ctx, buf := setupLogger(t)
		h := Handler(ctx)
		r := slog.Record{
			Time:    time.Now(),
			Message: "test message",
			Level:   slog.LevelInfo,
		}

		err := h.Handle(ctx, r)
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "INF test message")
	})

	t.Run("attrs", func(t *testing.T) {
		t.Parallel()

		ctx, buf := setupLogger(t)
		h := Handler(ctx)
		r := slog.Record{
			Time:    time.Now(),
			Message: "test message",
			Level:   slog.LevelInfo,
		}
		r.AddAttrs(slog.String("test", "test"))

		err := h.Handle(ctx, r)
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "INF test message test=test")
	})

}

func TestEnabled(t *testing.T) {
	h := handler{}

	t.Run("trace", func(t *testing.T) {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		assert.True(t, h.Enabled(t.Context(), slog.LevelDebug))
		assert.True(t, h.Enabled(t.Context(), slog.LevelInfo))
		assert.True(t, h.Enabled(t.Context(), slog.LevelWarn))
		assert.True(t, h.Enabled(t.Context(), slog.LevelError))
	})

	t.Run("debug", func(t *testing.T) {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		assert.True(t, h.Enabled(t.Context(), slog.LevelDebug))
		assert.True(t, h.Enabled(t.Context(), slog.LevelInfo))
		assert.True(t, h.Enabled(t.Context(), slog.LevelWarn))
		assert.True(t, h.Enabled(t.Context(), slog.LevelError))
	})

	t.Run("info", func(t *testing.T) {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		assert.False(t, h.Enabled(t.Context(), slog.LevelDebug))
		assert.True(t, h.Enabled(t.Context(), slog.LevelInfo))
		assert.True(t, h.Enabled(t.Context(), slog.LevelWarn))
		assert.True(t, h.Enabled(t.Context(), slog.LevelError))
	})

	t.Run("warn", func(t *testing.T) {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		assert.False(t, h.Enabled(t.Context(), slog.LevelDebug))
		assert.False(t, h.Enabled(t.Context(), slog.LevelInfo))
		assert.True(t, h.Enabled(t.Context(), slog.LevelWarn))
		assert.True(t, h.Enabled(t.Context(), slog.LevelError))
	})

	t.Run("error", func(t *testing.T) {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		assert.False(t, h.Enabled(t.Context(), slog.LevelDebug))
		assert.False(t, h.Enabled(t.Context(), slog.LevelInfo))
		assert.False(t, h.Enabled(t.Context(), slog.LevelWarn))
		assert.True(t, h.Enabled(t.Context(), slog.LevelError))
	})

	t.Run("fatal", func(t *testing.T) {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
		assert.False(t, h.Enabled(t.Context(), slog.LevelDebug))
		assert.False(t, h.Enabled(t.Context(), slog.LevelInfo))
		assert.False(t, h.Enabled(t.Context(), slog.LevelWarn))
		assert.False(t, h.Enabled(t.Context(), slog.LevelError))
	})
}

func TestSlogToZlogLevel(t *testing.T) {
	assert.Equal(t, zerolog.DebugLevel, slogToZlogLevel(slog.LevelDebug))
	assert.Equal(t, zerolog.InfoLevel, slogToZlogLevel(slog.LevelInfo))
	assert.Equal(t, zerolog.WarnLevel, slogToZlogLevel(slog.LevelWarn))
	assert.Equal(t, zerolog.ErrorLevel, slogToZlogLevel(slog.LevelError))
}

func TestSlogToZlogAttr(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(event, slog.String("test", "test"))
		event.Send()

		assert.Contains(t, buf.String(), "test=test")
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(event, slog.Bool("test", true))
		event.Send()

		assert.Contains(t, buf.String(), "test=true")
	})

	t.Run("duration", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(event, slog.Duration("test", time.Minute))
		event.Send()

		assert.Contains(t, buf.String(), "test=60000")
	})

	t.Run("float64", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(event, slog.Float64("test", 20.3))
		event.Send()

		assert.Contains(t, buf.String(), "test=20.3")
	})

	t.Run("time", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		testTime := time.Now()
		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(event, slog.Time("test", testTime))
		event.Send()

		assert.Contains(t, buf.String(), "test="+testTime.Format(time.RFC3339))
	})

	t.Run("uint64", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(event, slog.Uint64("test", uint64(21)))
		event.Send()

		assert.Contains(t, buf.String(), "test=21")
	})

	t.Run("group", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(
			event,
			slog.Group("test group",
				slog.Group("test subgroup",
					slog.Group("test subsubgroup",
						slog.String("test", "test"),
					),
				),
			),
		)
		event.Send()

		assert.Contains(t, buf.String(), `INF test group={"test subgroup":{"test subsubgroup":{"test":"test"}}}`)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(event, slog.Any("test", fmt.Errorf("bad things")))
		event.Send()

		assert.Contains(t, buf.String(), `INF error="bad things"`)
	})

	t.Run("any", func(t *testing.T) {
		t.Parallel()
		ctx, buf := setupLogger(t)

		event := log.Ctx(ctx).Info().Ctx(ctx)
		slogToZlogAttr(event, slog.Any("test", "test message"))
		event.Send()

		assert.Contains(t, buf.String(), `INF test="test message"`)
	})
}
