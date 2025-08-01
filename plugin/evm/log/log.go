// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package log

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"

	ethlog "github.com/ava-labs/libevm/log"
	"github.com/ava-labs/subnet-evm/log"
	"golang.org/x/exp/slog"
)

type Logger struct {
	ethlog.Logger

	logLevel *slog.LevelVar
}

// InitLogger initializes logger with alias and sets the log level and format with the original [os.StdErr] interface
// along with the context logger.
func InitLogger(alias string, level string, jsonFormat bool, writer io.Writer) (Logger, error) {
	logLevel := &slog.LevelVar{}

	var handler slog.Handler
	if jsonFormat {
		chainStr := fmt.Sprintf("%s Chain", alias)
		handler = log.JSONHandlerWithLevel(writer, logLevel)
		handler = &addContext{Handler: handler, logger: chainStr}
	} else {
		useColor := false
		chainStr := fmt.Sprintf("<%s Chain> ", alias)
		termHandler := log.NewTerminalHandlerWithLevel(writer, logLevel, useColor)
		termHandler.Prefix = func(r slog.Record) string {
			file, line := getSource(r)
			if file != "" {
				return fmt.Sprintf("%s%s:%d ", chainStr, file, line)
			}
			return chainStr
		}
		handler = termHandler
	}

	// Create handler
	c := Logger{
		Logger:   ethlog.NewLogger(handler),
		logLevel: logLevel,
	}

	if err := c.SetLogLevel(level); err != nil {
		return Logger{}, err
	}
	ethlog.SetDefault(c.Logger)
	return c, nil
}

// SetLogLevel sets the log level of initialized log handler.
func (l *Logger) SetLogLevel(level string) error {
	// Set log level
	logLevel, err := log.LvlFromString(level)
	if err != nil {
		return err
	}
	l.logLevel.Set(logLevel)
	return nil
}

// locationTrims are trimmed for display to avoid unwieldy log lines.
var locationTrims = []string{
	"subnet-evm/",
}

func trimPrefixes(s string) string {
	for _, prefix := range locationTrims {
		idx := strings.LastIndex(s, prefix)
		if idx >= 0 {
			s = s[idx+len(prefix):]
		}
	}
	return s
}

func getSource(r slog.Record) (string, int) {
	frames := runtime.CallersFrames([]uintptr{r.PC})
	frame, _ := frames.Next()
	return trimPrefixes(frame.File), frame.Line
}

type addContext struct {
	slog.Handler

	logger string
}

func (a *addContext) Handle(ctx context.Context, r slog.Record) error {
	r.Add(slog.String("logger", a.logger))
	file, line := getSource(r)
	if file != "" {
		r.Add(slog.String("caller", fmt.Sprintf("%s:%d", file, line)))
	}
	return a.Handler.Handle(ctx, r)
}
