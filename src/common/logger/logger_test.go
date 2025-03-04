package logger

import (
	"testing"
)

func TestSetLogSeverity(t *testing.T) {
	tests := []struct {
		name     string
		severity int
	}{
		{name: "SetFatal", severity: LogFatal},
		{name: "SetError", severity: LogError},
		{name: "SetInfo", severity: LogInfo},
		{name: "SetWarning", severity: LogWarning},
		{name: "SetDebug", severity: LogDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogSeverity(tt.severity)
			if lwInst.CurrentLogSeverity != tt.severity {
				t.Errorf("Expected severity %d, got %d", tt.severity, lwInst.CurrentLogSeverity)
			}
		})
	}
}

func TestStringifyEvent(t *testing.T) {
	tests := []struct {
		name  string
		param bool
	}{
		{name: "EnableStringify", param: true},
		{name: "DisableStringify", param: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StringifyEvent(tt.param)
			if lwInst.StringifyEvent != tt.param {
				t.Errorf("Expected StringifyEvent %v, got %v", tt.param, lwInst.StringifyEvent)
			}
		})
	}
}

func TestEventToEventId(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		event    interface{}
		expected map[string]interface{}
	}{
		{name: "BasicEvent", id: "event1", event: "testEvent", expected: map[string]interface{}{"event1": "testEvent"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := eventToEventId(tt.id, tt.event)
			if res, ok := result.(map[string]interface{}); ok {
				if res[tt.id] != tt.expected[tt.id] {
					t.Errorf("Expected %v, got %v", tt.expected, res)
				}
			} else {
				t.Errorf("Result was not a map")
			}
		})
	}
}

func TestSeverityToString(t *testing.T) {
	tests := []struct {
		name     string
		severity int
		expected string
	}{
		{name: "Fatal", severity: LogFatal, expected: logFatalString},
		{name: "Error", severity: LogError, expected: logErrorString},
		{name: "Info", severity: LogInfo, expected: logInfoString},
		{name: "Warning", severity: LogWarning, expected: logWarningString},
		{name: "Debug", severity: LogDebug, expected: logDebugString},
		{name: "Unknown", severity: 999, expected: logUnknownString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := severityToString(tt.severity)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBuildFilename(t *testing.T) {
	tests := []struct {
		name         string
		inputPath    string
		separator    string
		expectedPath string
		isSplit      bool
	}{
		{name: "UnixPath", inputPath: "/dir1/dir2/file.go", separator: "/", expectedPath: "dir2/file.go", isSplit: true},
		{name: "WindowsPath", inputPath: "C:\\dir1\\dir2\\file.go", separator: "\\", expectedPath: "dir2\\file.go", isSplit: true},
		{name: "NoSeparator", inputPath: "file.go", separator: "/", expectedPath: "file.go", isSplit: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := buildFilename(tt.inputPath, tt.separator)
			if result != tt.expectedPath || ok != tt.isSplit {
				t.Errorf("Expected (%s, %v), got (%s, %v)", tt.expectedPath, tt.isSplit, result, ok)
			}
		})
	}
}

func TestGetCallerStack(t *testing.T) {
	line, file, method := getCallerStack()
	if line == 0 || file == "" || method == "" {
		t.Errorf("Expected valid line, file, and method, got (%d, %s, %s)", line, file, method)
	}
}

func TestLogWithSeverity(t *testing.T) {
	tests := []struct {
		name     string
		severity int
		format   string
		args     []interface{}
	}{
		{name: "InfoLog", severity: LogInfo, format: "This is an %s", args: []interface{}{"info log"}},
		{name: "ErrorLog", severity: LogError, format: "This is an %s", args: []interface{}{"error log"}},
	}

	SetLogSeverity(LogDebug)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logWithSeverity(tt.severity, tt.format, tt.args...)
		})
	}
}
