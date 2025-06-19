package shared

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// SetupLogger sets up the log file and assigns it to Gin's log output
func SetupLogger() {
	// Configure lumberjack for log rotation
	logFile := &lumberjack.Logger{
		Filename:   "/app/runtime/app.log",
		MaxSize:    10,   // Max size in MB before rotation
		MaxBackups: 3,    // Max number of old log files to keep
		MaxAge:     28,   // Max days to retain old log files
		Compress:   true, // Compress the rotated log files
	}
	logFileError := &lumberjack.Logger{
		Filename:   "/app/runtime/error.log",
		MaxSize:    10,   // Max size in MB before rotation
		MaxBackups: 3,    // Max number of old log files to keep
		MaxAge:     28,   // Max days to retain old log files
		Compress:   true, // Compress the rotated log files
	}
	//
	//// Optionally log to both file and stdout
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(logFileError, os.Stderr)
	//gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	//	// your custom format
	//	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
	//		param.ClientIP,
	//		param.TimeStamp.Format(time.RFC1123),
	//		param.Method,
	//		param.Path,
	//		param.Request.Proto,
	//		param.StatusCode,
	//		param.Latency,
	//		param.Request.UserAgent(),
	//		param.ErrorMessage,
	//	)
	//})
}
