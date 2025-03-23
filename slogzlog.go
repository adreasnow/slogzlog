package slogzlog

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type slogzloghandler struct {
	slog.Handler
	ctx context.Context
}

// Creates a new slog handler, storing the stored context in the handler struct.
// This context should contain a zerolog.Logger that will be used by the handler
func NewSlogHandler(ctx context.Context) slogzloghandler {
	return slogzloghandler{ctx: ctx}
}

// Checks to see if the zerolog global log level is alllowed based on the incoming slog.Level
func (s slogzloghandler) Enabled(_ context.Context, level slog.Level) bool {
	var allowable []slog.Level
	switch zerolog.GlobalLevel() {
	case zerolog.TraceLevel:
		allowable = []slog.Level{slog.LevelError, slog.LevelWarn, slog.LevelInfo, slog.LevelDebug}
	case zerolog.DebugLevel:
		allowable = []slog.Level{slog.LevelError, slog.LevelWarn, slog.LevelInfo, slog.LevelDebug}
	case zerolog.InfoLevel:
		allowable = []slog.Level{slog.LevelError, slog.LevelWarn, slog.LevelInfo}
	case zerolog.WarnLevel:
		allowable = []slog.Level{slog.LevelError, slog.LevelWarn}
	case zerolog.ErrorLevel:
		allowable = []slog.Level{slog.LevelError}
	default:
		allowable = []slog.Level{}
	}

	return slices.Contains(allowable, level)
}

// Converts the slog.Record into a zerolog.Event and sends it using the logger
// that's stored in the context that was set when the handler was initilised
func (s slogzloghandler) Handle(ctx context.Context, r slog.Record) error {
	event := log.Ctx(s.ctx).
		WithLevel(slogToZlogLevel(r.Level)).
		Ctx(ctx)

	r.Attrs(func(attr slog.Attr) bool {
		slogToZlogAttr(event, attr)
		return true
	})

	event.Msg(r.Message)

	return nil
}

// Converts the slog.Level to the corresponding zerolog.Level
func slogToZlogLevel(l slog.Level) zerolog.Level {
	switch l {
	case slog.LevelDebug:
		return zerolog.DebugLevel
	case slog.LevelInfo:
		return zerolog.InfoLevel
	case slog.LevelWarn:
		return zerolog.WarnLevel
	case slog.LevelError:
		return zerolog.ErrorLevel
	}
	return zerolog.NoLevel
}

// Converts any slog.Attrs in the slog.Record into the appropriate zerolog type and
// attaches it to the *zerolog.Event. Works recursively for slog.Group, and will extract
// errors stored in a slog.Any() and set them as event.Err(err).
func slogToZlogAttr(event *zerolog.Event, attr slog.Attr) {
	if k, ok := attr.Value.Any().(error); ok {
		event.Err(k)
		return
	}

	switch attr.Value.Kind() {
	case slog.KindString:
		event.Str(attr.Key, attr.Value.String())
	case slog.KindBool:
		event.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		event.Dur(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		event.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		event.Int64(attr.Key, attr.Value.Int64())
	case slog.KindTime:
		event.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		event.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		dict := zerolog.Dict()
		for _, subAttr := range attr.Value.Group() {
			slogToZlogAttr(dict, subAttr)
		}
		event.Dict(attr.Key, dict)
	case slog.KindAny:
		fallthrough
	default:
		event.Str(attr.Key, fmt.Sprintf("%v", attr.Value.Any()))
	}
}
