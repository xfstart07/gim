// Author: xufei
// Date: 2019-07-19

package lg

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
)

var std = New(DebugLevel)

type LoggerPrintf interface {
	Printf(f string, v ...interface{})
}

type loggerPrintf struct {
	logger *zap.Logger
}

func NewPrintf(verbose string) LoggerPrintf {
	return &loggerPrintf{
		logger: New(verbose),
	}
}

func (l loggerPrintf) Printf(f string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(f, v...))
}

func Logger() *zap.Logger {
	return std
}

func SetLevel(lvl string) {
	if lvl == InfoLevel {
		std = New(lvl)
	}
}

func New(lvl string) *zap.Logger {
	var config zap.Config

	if lvl == DebugLevel {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.DisableStacktrace = false
	config.DisableCaller = false
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stdout"}

	lg, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}

	return lg
}

func Logf(f string, args ...interface{}) string {
	return fmt.Sprintf(f, args...)
}
