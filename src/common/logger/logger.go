package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogFatal represents a fatal level log severity, typically used for critical errors and program termination.
// LogError represents an error level log severity, used for non-critical errors requiring attention.
// LogWarning represents a warning level log severity, used for potentially harmful situations.
// LogInfo represents an informational level log severity, used for general informational messages.
// LogDebug represents a debug level log severity, used for detailed system or diagnostic information.
const (
	LogFatal   = -1
	LogError   = 0
	LogWarning = 1
	LogInfo    = 2
	LogDebug   = 3
)

// Constants representing log levels as string values.
const (
	logUnknownString = "UNKNOWN"
	logFatalString   = "FATAL"
	logErrorString   = "ERROR"
	logInfoString    = "INFO"
	logWarningString = "WARNING"
	logDebugString   = "DEBUG"
)

// RFC3339Milli defines a timestamp format with millisecond precision adhering to RFC 3339 standards.
const (
	RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
)

// logWrapper is used to manage log configurations including severity levels and options for event stringification.
// CurrentLogSeverity indicates the severity level of logs to be processed or displayed.
// StringifyEvent determines if log events should be converted to string format before logging.
type logWrapper struct {
	CurrentLogSeverity int
	StringifyEvent     bool
}

// lwInst is a pre-defined instance of logWrapper with default severity set to LogInfo and event stringification enabled.
var lwInst = &logWrapper{
	CurrentLogSeverity: LogInfo,
	StringifyEvent:     true,
}

// SetLogSeverity sets the current log severity level for the log wrapper instance.
func SetLogSeverity(severity int) {
	lwInst.CurrentLogSeverity = severity
}

// StringifyEvent sets the StringifyEvent property of the logWrapper instance to enable or disable event stringification.
func StringifyEvent(s bool) {
	lwInst.StringifyEvent = s
}

// jsonLogRow represents a single structured log entry in JSON format.
// It includes metadata such as timestamp, file name, line number, method name, log message, severity, and optional event data.
type jsonLogRow struct {
	Timestamp  string      `json:"timestamp"`
	FileName   string      `json:"file"`
	LineNumber int         `json:"line"`
	MethodName string      `json:"method"`
	Message    string      `json:"message"`
	Event      interface{} `json:"event,omitempty"`
	Severity   string      `json:"severity"`
}

//var defaultLocation, _ = time.LoadLocation("Europe/Rome")

// eventToEventId maps an event to an Id by wrapping it into a map with the given id as the key and the event as the value.
func eventToEventId(id string, event interface{}) interface{} {
	c := make(map[string]interface{})
	c[id] = event
	return c
}

// severityToString converts an integer severity level into a corresponding string representation of the log level.
func severityToString(severity int) string {
	var out string
	switch severity {
	case LogFatal:
		out = logFatalString
	case LogError:
		out = logErrorString
	case LogInfo:
		out = logInfoString
	case LogWarning:
		out = logWarningString
	case LogDebug:
		out = logDebugString
	default:
		out = logUnknownString
	}

	return out
}

// Fatal logs a fatal severity message, then exits the process with status code 255.
func Fatal(format string, a ...interface{}) {
	logWithSeverity(LogFatal, format, a...)
	os.Exit(255)
}

// Error logs a formatted error message with a predefined error severity level.
func Error(format string, a ...interface{}) {
	logWithSeverity(LogError, format, a...)
}

// Info logs a message with the severity level of LogInfo. The message format supports printf-style formatting.
func Info(format string, a ...interface{}) {
	logWithSeverity(LogInfo, format, a...)
}

// Warning logs a formatted message with a warning severity level.
func Warning(format string, a ...interface{}) {
	logWithSeverity(LogWarning, format, a...)
}

// Debug logs a formatted string with debug severity level. Accepts a format and optional arguments for formatting.
func Debug(format string, a ...interface{}) {
	logWithSeverity(LogDebug, format, a...)
}

// logWithSeverity logs a message with a given severity, including a formatted timestamp, file, line, and method information.
// It only logs messages if the severity is less than or equal to the current log severity configured in lwInst.
func logWithSeverity(severity int, format string, a ...interface{}) {
	t := time.Now().Local()
	timestamp := t.Format(RFC3339Milli)
	lineNumber, fileName, methodName := getCallerStack()
	formattedString := fmt.Sprintf(format, a...)
	if severity <= lwInst.CurrentLogSeverity {
		wholeRow := fmt.Sprintf("[%s][%s][%s][%d][%s] %s", timestamp, severityToString(severity), fileName, lineNumber, methodName, formattedString)
		log.Println(wholeRow)
	}
}

// buildFilename constructs a filename by combining the last two segments of the input string separated by the given separator.
// It returns the constructed filename and a boolean indicating whether the operation was successful.
func buildFilename(in string, sep string) (string, bool) {
	if fls := strings.Split(in, sep); len(fls) > 0 {
		if len(fls) > 1 {
			return fls[len(fls)-2] + string(os.PathSeparator) + fls[len(fls)-1], true
		}
	}
	return in, false
}

// getCallerStack retrieves the caller's line number, file name, and method name from the runtime stack.
func getCallerStack() (int, string, string) {
	lineNumber := 0
	fileName := ""
	methodName := ""
	pc := make([]uintptr, 10)
	runtime.Callers(3, pc)
	idx := 1
	if f := runtime.FuncForPC(pc[idx]); f != nil {
		fileName, lineNumber = f.FileLine(pc[idx])
		var ok bool
		if fileName, ok = buildFilename(fileName, "/"); !ok {
			fileName, _ = buildFilename(fileName, "\\")
		}
		methodName = f.Name()
		if pos := strings.LastIndex(methodName, "/"); pos >= 0 {
			if len(methodName) > pos+1 {
				pos++
			}
			methodName = methodName[pos:]
		}
	}
	return lineNumber, fileName, methodName
}
