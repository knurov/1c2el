package helper

import (
	"fmt"
	"path/filepath"

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

//Loger log structure
type Loger struct {
}

//Destroy closing opened files, databases and etc
func (log *Loger) Destroy() {
}

// NewLoger Create new loger
func NewLoger() (log *Loger) {
	log = new(Loger)
	log.Init("", logrus.ErrorLevel)
	return log
}

//Init Create new loger
func (log Loger) Init(logpath string, level logrus.Level) {
	if logpath != "" {
		logoutput := &lumberjack.Logger{
			Filename: filepath.Join(logpath, "1c2el.log"),
			MaxSize:  100, // megabytes
			// MaxBackups: 15,
			MaxAge:   15,   //days
			Compress: true, // disabled by default
		}
		logrus.SetOutput(logoutput)

	}
	logrus.SetLevel(level)
}

func (log Loger) msg(level logrus.Level, args ...interface{}) {
	var err *error
	var msg string

	switch len(args) {
	case 1:
		if args[0] == nil || args[0] == "" {
			return
		}
		msg = fmt.Sprint(args[0])
	case 2:
		if args[1] == nil {
			return
		}
		msg = fmt.Sprintf(args[0].(string), args[1])
	default:
		logrus.Fatalf("Incorerect param count for logMsg - %v", len(args))
	}

	switch level {
	case logrus.ErrorLevel:
		logrus.Error(msg)
	case logrus.FatalLevel:
		logrus.Fatal(msg)
	case logrus.DebugLevel:
		logrus.Debug(msg)
	case logrus.WarnLevel:
		logrus.Warn(msg)
	case logrus.InfoLevel:
		logrus.Info(msg)
	case logrus.TraceLevel:
		logrus.Trace(msg)
	default:
		logrus.Errorf("Unknov level!!! %v", fmt.Sprintf(msg, err))
	}

}

//Fatal Output fatal error to log
func (log Loger) Fatal(args ...interface{}) {
	log.msg(logrus.FatalLevel, args...)
}

//Error Otput error to log
func (log Loger) Error(args ...interface{}) {
	log.msg(logrus.ErrorLevel, args...)
}

//Warning Otput warning to log
func (log Loger) Warning(args ...interface{}) {
	log.msg(logrus.WarnLevel, args...)
}

//Debug Otput debug to log
func (log Loger) Debug(msg string, args ...interface{}) {
	logrus.Debugf(msg, args...)
}

//Info Otput info to log
func (log Loger) Info(msg string, args ...interface{}) {
	logrus.Infof(msg, args...)
}

//Trace Otput Trace to log
func (log Loger) Trace(msg string, args ...interface{}) {
	logrus.Tracef(msg, args...)
}
