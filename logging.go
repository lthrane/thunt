package main

import (
	"io"
	"log"
)

var (
	// Trace is a logger to provide debug trace logs
	Trace *log.Logger

	// Info is a logger to provide info logs
	Info *log.Logger

	// Warning is a logger to provide warning logs
	Warning *log.Logger

	// Error is a logger to provide error logs
	Error *log.Logger
)

// Init intializes the logging system
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
