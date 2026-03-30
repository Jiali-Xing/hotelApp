package config

import (
	"log"
	"runtime"
	"strings"
)

//var debug = flag.Bool("debug", false, "Enable debug logging")

func DebugLog(format string, v ...interface{}) {
	if Debug {
		// Get the caller function name
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		funcName := "unknown"
		if ok && details != nil {
			// Extract the function name
			fullFuncName := details.Name()
			// Remove the package path
			parts := strings.Split(fullFuncName, "/")
			funcName = parts[len(parts)-1]
		}
		log.Printf("DEBUG [%s]: "+format, append([]interface{}{funcName}, v...)...)
	}
}
