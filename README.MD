# Slog to Zerolog bridge

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/adreasnow/slogzlog) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/adreasnow/slogzlog/main/LICENSE) [![Build Status](https://github.com/adreasnow/slogzlog/actions/workflows/test-tag.yaml/badge.svg?branch=main)](https://github.com/adreasnow/slogzlog/actions/workflows/test-tag.yaml) [![Go Coverage](https://github.com/adreasnow/slogzlog/wiki/coverage.svg)](https://raw.githack.com/wiki/adreasnow/slogzlog/coverage.html)

The slogzlog package provides a handler for slog that sends incoming logs to zerolog. The module was created as I love the [github.com/go-fuego/fuego](https://github.com/go-fuego/fuego) router, however it only supports passing in a `slog.Handler` to the logging field, and I much prefer using [github.com/rs/zerolog](https://github.com/rs/zerolog).

## Usage

```go
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
	handler := New(ctx)
```
