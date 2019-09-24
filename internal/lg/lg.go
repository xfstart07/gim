// Author: xufei
// Date: 2019-07-19

package lg

import (
	"flag"
	"fmt"
	"log"

	"go.uber.org/zap"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
)

var (
	stdLevel string
	std      *zap.Logger
)

func init() {
	// FIXME：flag 是隐藏在不同的包中，还是统一放在 main 中好呢？
	flag.StringVar(&stdLevel, "log_level", DebugLevel, "set server log level")

	std = New(stdLevel)
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
