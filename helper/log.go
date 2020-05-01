package helper

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func logMsg(level log.Level, args ...interface{}) {
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
		log.Fatalf("Incorerect param count for logMsg - %v", len(args))
	}

	// for _, arg := range args {
	// 	// if reflect.TypeOf()
	// 	switch argType := arg.(type) {
	// 	// case interface{}:
	// 	case string:
	// 		msg = arg.(string)
	// 	case error:
	// 		err = arg.(*error)
	// 	default:
	// 		log.Fatalf("Uknov type %v", argType)
	// 	}
	// }
	// if msg == "" {
	// 	msg = "v%"
	// }

	switch level {
	case log.ErrorLevel:
		log.Error(msg)
	case log.FatalLevel:
		log.Fatal(msg)
	default:
		log.Errorf("Unknov level!!! %v", fmt.Sprintf(msg, err))
	}

}

//LogFatal Output fatal error to log
func LogFatal(args ...interface{}) {
	logMsg(log.FatalLevel, args...)
}

//LogError Otput error to log
func LogError(args ...interface{}) {
	logMsg(log.ErrorLevel, args...)
}
