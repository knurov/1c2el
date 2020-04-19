package main

import (
	"flag"
	"fmt"
	"path/filepath"

	// "log"
	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// import args

func main() {
	var msg string
	configName := flag.String("config", "./config.yaml", "Set config file path")
	logpath := flag.String("log", "", "Set log path")
	flag.StringVar(&msg, "message", "hello!", "just message")
	flag.Parse()
	fmt.Printf("Begin! %v\n", *configName)
	fmt.Printf("%s end\n", msg)

	if *logpath != "" {
		logoutput := &lumberjack.Logger{
			Filename: filepath.Join(*logpath, "1c2el.log"),
			MaxSize:  100, // megabytes
			// MaxBackups: 15,
			MaxAge:   15,   //days
			Compress: true, // disabled by default
		}
		log.SetOutput(logoutput)

	}
	log.SetLevel(log.TraceLevel)
	log.Errorf("Err msq %v", msg)
	log.Tracef("Trace val %v", log.TraceLevel)
	log.Println("bue!")
}
