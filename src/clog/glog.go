package clog

import (
	"fmt"
	"github.com/golang/glog"
)

func init() {
	LoggerRegistry["glog"] = GLogger{}
}

type GLogger struct{}

func (GLogger) Infof(format string, args ...interface{}) {
	glog.InfoDepth(3, fmt.Sprintf(format, args...))
}
func (GLogger) Warningf(format string, args ...interface{}) {
	glog.WarningDepth(3, fmt.Sprintf(format, args...))
}
func (GLogger) Errorf(format string, args ...interface{}) {
	glog.ErrorDepth(3, fmt.Sprintf(format, args...))
}
func (GLogger) Fatalf(format string, args ...interface{}) {
	glog.FatalDepth(3, fmt.Sprintf(format, args...))
}

