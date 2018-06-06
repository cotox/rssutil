// Copyright 2018 cotox. All rights reserved.
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

// TODO: add color feature

package rssutil

import (
	"fmt"
	"log"
	"os"
)

var LogLevel = Lerror

var traceLogger = log.New(os.Stderr, "", log.LstdFlags+log.Lshortfile)
var debugLogger = log.New(os.Stderr, "", log.LstdFlags+log.Lshortfile)
var infoLogger = log.New(os.Stderr, "", log.LstdFlags+log.Lshortfile)
var warnLogger = log.New(os.Stderr, "", log.LstdFlags+log.Lshortfile)
var errLogger = log.New(os.Stderr, "", log.LstdFlags+log.Lshortfile)

const (
	// Ltrace indicates log trace level info
	Ltrace = iota
	// Ldebug indicates log debug level info
	Ldebug
	// Linfo indicates log info level info
	Linfo
	// Lwarning indicates log warning level info
	Lwarning
	// Lerror indicates log error level info
	Lerror
)

func logTracef(format string, v ...interface{}) {
	if LogLevel == Ltrace {
		traceLogger.Output(2, fmt.Sprintf("[TRACE] "+format, v...))
	}
}
func logTrace(v ...interface{}) {
	if LogLevel == Ltrace {
		var v2 []interface{}
		v2 = append(v2, "[TRACE] ")
		v2 = append(v2, v...)
		traceLogger.Output(2, fmt.Sprint(v2...))
	}
}
func logTraceln(v ...interface{}) {
	if LogLevel == Ltrace {
		var v2 []interface{}
		v2 = append(v2, "[TRACE] ")
		v2 = append(v2, v...)
		traceLogger.Output(2, fmt.Sprintln(v2...))
	}
}

func logDebugf(format string, v ...interface{}) {
	if LogLevel <= Ldebug {
		debugLogger.Output(2, fmt.Sprintf("[DEBUG] "+format, v...))
	}
}
func logDebug(v ...interface{}) {
	if LogLevel <= Ldebug {
		var v2 []interface{}
		v2 = append(v2, "[DEBUG] ")
		v2 = append(v2, v...)
		debugLogger.Output(2, fmt.Sprint(v2...))
	}
}
func logDebugln(v ...interface{}) {
	if LogLevel <= Ldebug {
		var v2 []interface{}
		v2 = append(v2, "[DEBUG] ")
		v2 = append(v2, v...)
		debugLogger.Output(2, fmt.Sprintln(v2...))
	}
}

func logInfof(format string, v ...interface{}) {
	if LogLevel <= Ldebug {
		infoLogger.Output(2, fmt.Sprintf("[DEBUG] "+format, v...))
	}
}
func logInfo(v ...interface{}) {
	if LogLevel <= Ldebug {
		var v2 []interface{}
		v2 = append(v2, "[DEBUG] ")
		v2 = append(v2, v...)
		infoLogger.Output(2, fmt.Sprint(v2...))
	}
}
func logInfoln(v ...interface{}) {
	if LogLevel <= Ldebug {
		var v2 []interface{}
		v2 = append(v2, "[DEBUG] ")
		v2 = append(v2, v...)
		infoLogger.Output(2, fmt.Sprintln(v2...))
	}
}

func logWarnf(format string, v ...interface{}) {
	if LogLevel <= Lwarning {
		warnLogger.Output(2, fmt.Sprintf("[WARN] "+format, v...))
	}
}
func logWarn(v ...interface{}) {
	if LogLevel <= Lwarning {
		var v2 []interface{}
		v2 = append(v2, "[WARN] ")
		v2 = append(v2, v...)
		warnLogger.Output(2, fmt.Sprint(v2...))
	}
}
func logWarnln(v ...interface{}) {
	if LogLevel <= Lwarning {
		var v2 []interface{}
		v2 = append(v2, "[WARN] ")
		v2 = append(v2, v...)
		warnLogger.Output(2, fmt.Sprintln(v2...))
	}
}

func logErrf(format string, v ...interface{}) {
	if LogLevel <= Lerror {
		errLogger.Output(2, fmt.Sprintf("[ERROR] "+format, v...))
	}
}
func logErr(v ...interface{}) {
	if LogLevel <= Lerror {
		var v2 []interface{}
		v2 = append(v2, "[ERROR] ")
		v2 = append(v2, v...)
		errLogger.Output(2, fmt.Sprint(v2...))
	}
}
func logErrln(v ...interface{}) {
	if LogLevel <= Lerror {
		var v2 []interface{}
		v2 = append(v2, "[ERROR] ")
		v2 = append(v2, v...)
		errLogger.Output(2, fmt.Sprintln(v2...))
	}
}
