package slogzlog

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestExampleUsage(t *testing.T) {
	// Set the global log level as desired. Slog's slowest level is debug, so
	// zerolog.TraceLevel will match this behaviour. Since slog has no fatal or
	// panic levels, these will effectively disable the slogzlog bridge
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Create a new zerolog.Logger and store it in a context. This logger also
	// outputs to a buffer for testing, as well as to *testing.T
	buf := new(bytes.Buffer)
	ctx := log.
		Output(
			io.MultiWriter(
				zerolog.TestWriter{T: t},
				zerolog.ConsoleWriter{Out: os.Stdout},
				zerolog.ConsoleWriter{Out: buf, NoColor: true},
			),
		).
		WithContext(t.Context())

	// Create te slogzlog handler with the previously created logger. The context
	// containing the logger will be stored in the handler
	handler := NewSlogHandler(ctx)

	// Use the handler as desired

	{ // Log an error
		startTime := time.Now()
		r := slog.Record{
			Time:    time.Now(),
			Message: "an error occurred",
			Level:   slog.LevelError,
		}
		r.AddAttrs(
			slog.Any("error", fmt.Errorf("something bad happened...")),
			slog.Duration("time_to_error", time.Since(startTime)),
		)
		handler.Handle(ctx, r)
	}

	// Check the message that got logged to zerolog

	assert.Contains(t, buf.String(), `ERR an error occurred error="something bad happened..." time_to_error=0.`)
}
