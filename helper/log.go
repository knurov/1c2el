package helper

import (
	log "github.com/sirupsen/logrus"
)

//LogFatal Output fatal error to log
func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//LogError Otput error to log
func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
