package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	// "log"
	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

func main() {
	var msg string
	configName := *(flag.String("config", "./config.yaml", "Set config file path"))
	logpath := *(flag.String("log", "", "Set log path"))
	flag.StringVar(&msg, "message", "hello!", "just message")
	flag.Parse()
	fmt.Printf("Begin! %v\n", configName)
	fmt.Printf("%s end\n", msg)

	type Config struct {
		Database struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"database"`
	}

	var config Config

	if configName != "" {
		configFile, confErr := ioutil.ReadFile(configName)
		if confErr != nil {
			log.Fatal(confErr)
		}
		yamlErr := yaml.Unmarshal(configFile, &config)
		if yamlErr != nil {
			log.Fatal(yamlErr)
		}
		fmt.Printf("Database.Host: %#v\n", config.Database.Host)

	}

	if logpath != "" {
		logoutput := &lumberjack.Logger{
			Filename: filepath.Join(logpath, "1c2el.log"),
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
