package logging

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
	hclog "github.com/hashicorp/go-hclog"
)

var statsPrefix = "service.facedetection."

// Logger defines an interface for common logging operations
type Logger interface {
	Log() hclog.Logger
	ServiceStart(address, port string)
	HealthHandlerCalled() Finished
}

// Finished defines a function to be returned by logging methods which contain timers
type Finished func()

// New creates a new logger with the given name and points it at a statsd server
func New(name, statsDServer, logLevel string) (Logger, error) {
	o := hclog.DefaultOptions
	o.Name = name
	o.Level = hclog.LevelFromString(logLevel)
	l := hclog.New(o)

	c, err := statsd.New(statsDServer)

	if err != nil {
		return nil, err
	}

	return &LoggerImpl{l, c}, nil
}

// LoggerImpl is a concrete implementation for the logger function
type LoggerImpl struct {
	l hclog.Logger
	s *statsd.Client
}

// Log returns the underlying logger
func (l *LoggerImpl) Log() hclog.Logger {
	return l.l
}

// ServiceStart logs information about the service start
func (l *LoggerImpl) ServiceStart(address, port string) {
	l.s.Incr(statsPrefix+"started", nil, 1)
	l.l.Info("Service started", "address", address, "port", port)
}

// HealthHandlerCalled logs information when the health handler is called, the returned function
// must be called once work has completed
func (l *LoggerImpl) HealthHandlerCalled() Finished {
	st := time.Now()

	return func() {
		l.s.Timing(statsPrefix+"health.called", time.Now().Sub(st), nil, 1)
		l.l.Debug("Health handler called")
	}
}
