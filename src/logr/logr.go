package logr

import (
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	FieldComponent = "component"
)

func Get() *logrus.Entry {
	logger := logrus.New()
	logger.Out = os.Stdout
	_, callerPath, _, ok := runtime.Caller(1)
	if ok {
		dir, file := path.Split(callerPath)
		return logger.WithField(FieldComponent, path.Join(path.Base(dir), file))
	} else {
		return logger.WithFields(logrus.Fields{})
	}
}
