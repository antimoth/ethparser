package log

import (
	"fmt"
	"path"
	"runtime"
)

func GetRuntimeInfo(skip int) string {
	pc, fileP, line, ok := runtime.Caller(skip)
	if !ok {
		return "Get runtime caller error!"
	}
	return fmt.Sprintf("%s@%s:%d", runtime.FuncForPC(pc).Name(), path.Base(fileP), line)
}

func toLogLevel(lvl string) logLevel {
	switch lvl {
	case "debug", "DEBUG":
		return DEBUG
	case "info", "INFO":
		return INFO
	case "warn", "WARN":
		return WARN
	case "error", "ERROR":
		return ERROR
	case "critical", "CRITICAL":
		return CRITICAL
	default:
		return DEBUG
	}
}
