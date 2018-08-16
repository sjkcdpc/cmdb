package logrus_stack

import (
	"strconv"

	"github.com/heralight/logrus_mate"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus_mate.RegisterHook("stackhook", NewStackHook)
}

func getSubLevels(levelIndex int) []logrus.Level {
	return logrus.AllLevels[0:levelIndex]
}

func NewStackHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	// Set levels to CallerLevels for which "caller" value may be set,
	// providing a single frame of stack.
	var callerLevels []logrus.Level
	if level, ok := options["caller-level"]; ok {
		if i, err := strconv.Atoi(level.(string)); err == nil {
			callerLevels = getSubLevels(i)
		}
	} else {
		callerLevels = logrus.AllLevels
	}

	var stackLevels []logrus.Level
	if level, ok := options["stack-level"]; ok {
		if i, err := strconv.Atoi(level.(string)); err == nil {
			stackLevels = getSubLevels(i)
		}
	} else {
		stackLevels = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
	}
	return LogrusStackHook{
		CallerLevels: callerLevels,
		StackLevels:  stackLevels,
	}, nil

}
