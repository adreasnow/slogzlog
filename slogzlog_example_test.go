package slogzlog

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ExampleNew() {
	// Set the global log level as desired. Slog's slowest level is debug, so
	// zerolog.TraceLevel will match this behaviour. Since slog has no fatal or
	// panic levels, these will effectively disable the slogzlog bridge
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Create a new zerolog.Logger
	buf := new(bytes.Buffer)
	zedLogger := log.
		Output(
			io.MultiWriter(
				zerolog.ConsoleWriter{Out: buf, NoColor: true},
			),
		)

	// Create the slogzlog handler with the previously created logger
	handler := New(&zedLogger)

	// Use the handler as desired
	logger := slog.New(handler)
	logger.Error("an error occurred",
		slog.Any("error", fmt.Errorf("something bad happened")),
		slog.Duration("time_to_error", time.Millisecond*125),
	)

	// Checks that the slog request has been processed by zerolog (through slogzlog)
	fmt.Println(strings.Join(strings.Split(buf.String(), " ")[1:], " "))

	// Output:
	// ERR an error occurred error="something bad happened" time_to_error=125
}
