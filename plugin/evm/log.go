// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

const (
	errorKey   = "LOG15_ERROR"
	timeFormat = "2006-01-02T15:04:05.000000-0700"
)

type SubnetEVMLogger struct {
	log.Handler
}

// InitLogger initializes logger with alias and sets the log level and format with the original [os.StdErr] interface
// along with the context logger.
func InitLogger(alias string, level string, jsonFormat bool, writer io.Writer) (SubnetEVMLogger, error) {
	// logFormat := SubnetEVMTermFormat(alias)
	logFormat := LogfmtFormat()
	if jsonFormat {
		logFormat = SubnetEVMJSONFormat(alias)
	}

	// Create handler
	logHandler := log.StreamHandler(writer, logFormat)
	logHandler = log.CallerFileHandler(logHandler)
	logHandler = HubbleTypeHandler(logHandler)
	logHandler = HubbleErrorHandler(logHandler)
	c := SubnetEVMLogger{Handler: logHandler}

	if err := c.SetLogLevel(level); err != nil {
		return SubnetEVMLogger{}, err
	}
	return c, nil
}

// SetLogLevel sets the log level of initialized log handler.
func (c *SubnetEVMLogger) SetLogLevel(level string) error {
	// Set log level
	logLevel, err := log.LvlFromString(level)
	if err != nil {
		return err
	}
	log.Root().SetHandler(log.LvlFilterHandler(logLevel, c))
	return nil
}

func SubnetEVMTermFormat(alias string) log.Format {
	prefix := fmt.Sprintf("<%s Chain>", alias)
	return log.FormatFunc(func(r *log.Record) []byte {
		location := fmt.Sprintf("%+v", r.Call)
		newMsg := fmt.Sprintf("%s %s: %s", prefix, location, r.Msg)
		r.Msg = newMsg
		return log.TerminalFormat(false).Format(r)
	})
}

func SubnetEVMJSONFormat(alias string) log.Format {
	prefix := fmt.Sprintf("%s Chain", alias)
	return log.FormatFunc(func(r *log.Record) []byte {
		props := make(map[string]interface{}, 5+len(r.Ctx)/2)
		props["timestamp"] = r.Time
		props["level"] = r.Lvl.String()
		props[r.KeyNames.Msg] = r.Msg
		props["logger"] = prefix
		props["caller"] = fmt.Sprintf("%+v", r.Call)
		for i := 0; i < len(r.Ctx); i += 2 {
			k, ok := r.Ctx[i].(string)
			if !ok {
				props[errorKey] = fmt.Sprintf("%+v is not a string key", r.Ctx[i])
			} else {
				// The number of arguments is normalized from the geth logger
				// to ensure that this will not cause an index out of bounds error
				props[k] = formatJSONValue(r.Ctx[i+1])
			}
		}

		b, err := json.Marshal(props)
		if err != nil {
			b, _ = json.Marshal(map[string]string{
				errorKey: err.Error(),
			})
			return b
		}

		b = append(b, '\n')
		return b
	})
}

func HubbleTypeHandler(h log.Handler) log.Handler {
	return log.FuncHandler(func(r *log.Record) error {
		var logType string
		// works for evm/limit_order.go, evm/orderbook/*.go, precompile/contracts/*
		if containsAnySubstr(r.Call.Frame().File, []string{"orderbook", "limit_order", "contracts"}) {
			logType = "hubble"
		} else {
			logType = "system"
		}
		// it's also possible to add type=hubble in logs originating from other files
		// by setting logtype=hubble and checking for it in this function by iterating through r.Ctx
		r.Ctx = append(r.Ctx, "type", logType)
		return h.Log(r)
	})
}

func HubbleErrorHandler(h log.Handler) log.Handler {
	// sets lvl=error when key name is "err" and value is not nil
	return log.FuncHandler(func(r *log.Record) error {
		for i := 0; i < len(r.Ctx); i += 2 {
			if r.Ctx[i] == "err" && r.Ctx[i+1] != nil {
				r.Lvl = log.LvlError
			}
		}
		return h.Log(r)
	})
}

func ErrorOnlyHandler(h log.Handler) log.Handler {
	// ignores all logs except lvl=error
	return log.FuncHandler(func(r *log.Record) error {
		if r.Lvl == log.LvlError {
			return h.Log(r)
		}
		return nil
	})
}

func formatJSONValue(value interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr && v.IsNil() {
				result = "nil"
			} else {
				panic(err)
			}
		}
	}()

	switch v := value.(type) {
	case time.Time:
		return v.Format(timeFormat)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

// containsAnySubstr checks if the string contains any of the specified substrings
func containsAnySubstr(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
