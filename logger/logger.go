package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
)

func Init() {

	// common flags: date, time, and short file:line
	flags := log.Ldate | log.Ltime | log.Lshortfile

	// Info logs to stdout
	InfoLogger = log.New(os.Stdout,
		"INFO: ",
		flags,
	)

	// DEBUG logs to stdout (or discard in production)
	DebugLogger = log.New(os.Stdout,
		"DEBUG: ",
		flags,
	)

	// To disable debug logs, in production we could do is
	// DebugLogger.SetOutput(io.Discard)

	ErrorLogger = log.New(os.Stderr,
		"Error: ",
		flags,
	)
}
