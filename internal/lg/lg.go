// Author: xufei
// Date: 2019-07-19

package lg

import (
	"flag"
	"log"

	"go.uber.org/zap/zapcore"

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
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// FIXME: 采用 pkg/errors 的 WithStack 来打印调用栈
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
