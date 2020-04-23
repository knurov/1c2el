package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	// "log"
	"github.com/antchfx/xmlquery"
	yamlconvert "github.com/ghodss/yaml"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	gojsonschema "github.com/xeipuuv/gojsonschema"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"launchpad.net/xmlpath"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func xmlparseOld() {
	file, err := os.Open("exaples/data16.04.2020 10_15_59.xml")
	fatal(err)
	doc, err := xmlquery.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	list := xmlquery.Find(doc, "//s")
	fmt.Println(list)

}

func xmlparseXmlpath() {
	file, err := os.Open("exaples/data16.04.2020 10_15_59.xml")
	defer file.Close()
	fatal(err)

	doc, err := xmlpath.Parse(file)
	fatal(err)

	path := xmlpath.MustCompile("/ТР/ОписаниеТрансформатора")
	// if list, ok := path.String(doc); ok {
	// 	fmt.Println(list)
	// }
	// fatal(err)
	// range item := path.Iter(doc) {
	items := path.Iter(doc)
	for items.Next() {
		value := items.Node()
		fmt.Println(value.String())
	}
}

func xmlparse() {

	type TRParams struct {
		Name string `xml:"Название,attr"`
	}

	type TRDescription struct {
		Number string   `xml:"НомерТрансформатора,attr"`
		Order  string   `xml:"НомерЗаказаНаПроизводство,attr"`
		Series string   `xml:"НомерСерии,attr"`
		Params TRParams `xml:"ПараметрыТрансформатора"`
	}

	type TR struct {
		Transfonmer xml.Name        `xml:"ТР"`
		Description []TRDescription `xml:"ОписаниеТрансформатора"`
	}
	file, err := os.Open("exaples/data16.04.2020 10_15_59.xml")
	defer file.Close()
	fatal(err)
	result := TR{}
	xmlData, err := ioutil.ReadAll(file)
	fatal(err)
	err = xml.Unmarshal(xmlData, &result)
	fmt.Print(result)
}

func db() {
	// postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]

	// postgresql://
	// postgresql://localhost
	// postgresql://localhost:5433
	// postgresql://localhost/mydb
	// postgresql://user@localhost
	// postgresql://user:secret@localhost
	// postgresql://other@localhost/otherdb?connect_timeout=10&application_name=myapp
	// postgresql://host1:123,host2:456/somedb?target_session_attrs=any&application_name=myapp

	db, err := pgx.Connect(context.Background(), "postgresql://electrolab:electrolab@localhost/electrolab")
	fatal(err)

	res := db.Ping(context.Background())
	fmt.Println(res)
	defer db.Close(context.Background())
}

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
	// https://json-schema.org/learn/miscellaneous-examples.html
	schema := gojsonschema.NewStringLoader(`
	{
		"required": [ "database" ],
		"properties": {
			"database": {
				"type": "object",
				"required": [ "host", "port" ],
				"host": {
					"type": "string"
				},
				"port": {
					"type": "int"
				}
			}
		}
	} 
	`)

	var config Config

	if configName != "" {
		configFile, err := ioutil.ReadFile(configName)
		fatal(err)
		configJSON, err := yamlconvert.YAMLToJSON(configFile)
		fatal(err)
		validationResult, err := gojsonschema.Validate(schema, gojsonschema.NewBytesLoader(configJSON))
		fatal(err)

		if !validationResult.Valid() {
			log.Fatal(validationResult.Errors())
		}

		err = yaml.Unmarshal(configFile, &config)
		fatal(err)
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
	db()
	xmlparse()
	log.Println("bue!")
}
