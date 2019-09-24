// Author: xufei
// Date: 2019-09-24 17:01

package lg

import (
	"fmt"
	"go.uber.org/zap"
)

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
