package logger

import (
	"log"
	"os"
)

const (
	InfoLevel  string = "info"
	WarnLevel  string = "warn"
	DebugLevel string = "debug"
	ErrorLevel string = "error"
	FatalLevel string = "fatal"
)

type Provider interface {
	Infow(msg string, keysAndValues ...interface{})
	Infof(format string, a ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Errorf(format string, a ...interface{})
	// Fatal logs a fatal error and exit with status 1.
	Fatal(msg string, err error)
	// Close the logger safely.
	Close()
}

type provider struct{}

// This provider is for a local mock logger provider. The real development env
// will implement above Provider interface and do necessary work like sending logs to
// cloud storage.
func NewProvider() Provider {
	return &provider{}
}

func (m *provider) Infow(msg string, keysAndValues ...interface{}) {
	log.Printf("info is: %s and kvs are: %+v\n\n", msg, keysAndValues)
}

func (m *provider) Infof(format string, a ...interface{}) {
	log.Printf("info is: %s and kvs are: %+v\n\n", format, a)
}

func (m *provider) Warnw(msg string, keysAndValues ...interface{}) {
	log.Printf("warn is: %s and kvs are: %+v\n\n", msg, keysAndValues)
}

func (m *provider) Debugw(msg string, keysAndValues ...interface{}) {
	log.Printf("debug is: %s and kvs are: %+v\n\n", msg, keysAndValues)
}

func (m *provider) Errorw(msg string, keysAndValues ...interface{}) {
	log.Printf("error is: %s and kvs are: %+v\n\n", msg, keysAndValues)
}

func (m *provider) Fatal(msg string, err error) {
	log.Printf("fatal is: %s and error is: %+v\n\n", msg, err)
	os.Exit(1)
}

func (m *provider) Errorf(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func (m *provider) Close() {
	// do nothing
}
