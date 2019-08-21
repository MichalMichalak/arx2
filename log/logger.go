package log

import (
	"fmt"
	"time"
)

type Logger interface {
	Log(severity Severity, v ...interface{})
	Logf(severity Severity, format string, params ...interface{})
	Debug(v ...interface{})
	Debugf(format string, params ...interface{})
	Info(v ...interface{})
	Infof(format string, params ...interface{})
	Warn(v ...interface{})
	Warnf(format string, params ...interface{})
	Error(v ...interface{})
	Errorf(format string, params ...interface{})
	Flush()
}

func NewServiceLogger(level Severity) ServiceLogger {
	return ServiceLogger{levels: level.ThisAndAbove()}
}

type ServiceLogger struct {
	levels map[Severity]struct{}
}

func (sl ServiceLogger) Log(severity Severity, v ...interface{}) {
	// Temporary implementation
	if _, ok := sl.levels[severity]; !ok {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("[" + now + "][" + string(severity) + "] ")
	fmt.Print(v...)
	defer fmt.Print("\n")

}

func (sl ServiceLogger) Logf(severity Severity, format string, params ...interface{}) {
	// Temporary implementation
	if _, ok := sl.levels[severity]; !ok {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("["+now+"]["+string(severity)+"] "+format, params...)
	defer fmt.Print("\n")
}

func (sl ServiceLogger) Debug(v ...interface{}) {
	if _, ok := sl.levels[SeverityDebug]; ok {
		sl.Log(SeverityDebug, v...)
	}
}

func (sl ServiceLogger) Debugf(format string, params ...interface{}) {
	if _, ok := sl.levels[SeverityDebug]; ok {
		sl.Logf(SeverityDebug, format, params...)
	}
}

func (sl ServiceLogger) Info(v ...interface{}) {
	if _, ok := sl.levels[SeverityInfo]; ok {
		sl.Log(SeverityInfo, v...)
	}
}

func (sl ServiceLogger) Infof(format string, params ...interface{}) {
	if _, ok := sl.levels[SeverityInfo]; ok {
		sl.Logf(SeverityInfo, format, params...)
	}
}

func (sl ServiceLogger) Warn(v ...interface{}) {
	if _, ok := sl.levels[SeverityWarning]; ok {
		sl.Log(SeverityWarning, v...)
	}
}

func (sl ServiceLogger) Warnf(format string, params ...interface{}) {
	if _, ok := sl.levels[SeverityWarning]; ok {
		sl.Logf(SeverityWarning, format, params...)
	}
}

func (sl ServiceLogger) Error(v ...interface{}) {
	if _, ok := sl.levels[SeverityError]; ok {
		sl.Log(SeverityError, v...)
	}
}

func (sl ServiceLogger) Errorf(format string, params ...interface{}) {
	if _, ok := sl.levels[SeverityError]; ok {
		sl.Logf(SeverityError, format, params...)
	}
}

func (sl ServiceLogger) Flush() {
	// Not used
}
