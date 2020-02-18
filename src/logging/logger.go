package logging

import (
	"os"

	"github.com/teambition/gear"
	gearLogging "github.com/teambition/gear/logging"
)

// Logger ...
var Logger = gearLogging.New(os.Stdout)

func init() {
	Logger.SetJSONLog()
}

// SetLevel ...
func SetLevel(level string) {
	l, err := gearLogging.ParseLevel(level)
	if err == nil {
		Logger.SetLevel(l)
	} else {
		Logger.Err(err)
	}
}

// Err produce a "Error" log with the default logger
func Err(v interface{}) {
	Logger.Err(v)
}

// Warning produce a "Warning" log with the default logger
func Warning(v interface{}) {
	Logger.Warning(v)
}

// Info produce a "Informational" log with the default logger
func Info(v interface{}) {
	Logger.Info(v)
}

// Debug produce a "Debug" log with the default logger
func Debug(v interface{}) {
	Logger.Debug(v)
}

// Debugf produce a "Debug" log in the manner of fmt.Printf with the default logger
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Panic produce a "Emergency" log with the default logger and then calls panic with the message
func Panic(v interface{}) {
	Logger.Panic(v)
}

// FromCtx retrieve the Log instance for the default logger.
func FromCtx(ctx *gear.Context) gearLogging.Log {
	return Logger.FromCtx(ctx)
}
