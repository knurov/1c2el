package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4"
	"knurov.ru/el/1c2el/config"
	"knurov.ru/el/1c2el/helper"
)

// func parseTLO10(fullName string)
// https://yourbasic.org/golang/regexp-cheat-sheet/
func parseFullName(log *helper.Loger, fullName string) {
	isTlo10, err := regexp.Compile("^ТЛО-10")
	// 'ТЛО-10_М1ACE-0.2SFS7/0.5FS7/10P10-10/10/40-150(300)-150(300)-300/5 У2 б 31.5кА'
	// tlo10, err := regexp.Compile("^(ТЛО-10)_(.+)-(.+)/(.+)/(.+)/(.+)/(.+)/(.+) (.+) (.+) (.+)")
	tlo10, err := regexp.Compile(`(?P<short>.+?)_(?P<prop>.+?)-`)
	log.Fatal(err)
	if isTlo10.MatchString(fullName) {
		result := tlo10.FindStringSubmatch(fullName)
		fmt.Println(tlo10.SubexpNames())
		for _, item := range result {
			fmt.Println(item)
		}
		fmt.Println(result)
	}
}

func xmlparse(log *helper.Loger, fileName string) {

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
	file, err := os.Open(fileName)
	defer file.Close()
	log.Fatal(err)
	result := TR{}
	xmlData, err := ioutil.ReadAll(file)
	log.Fatal(err)
	err = xml.Unmarshal(xmlData, &result)

	for _, item := range result.Description {
		// fmt.Printf("%v - Serial number %v\n", item.Params.Name, item.Number)
		parseFullName(log, item.Params.Name)
	}
}

func getFiles(log *helper.Loger, path string) {
	files, err := ioutil.ReadDir(path)
	log.Fatal(err)
	for _, item := range files {
		if !item.IsDir() {
			xmlparse(log, filepath.Join(path, item.Name()))
		}

	}
}

func db(log *helper.Loger) {
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
	defer db.Close(context.Background())
	log.Fatal(err)

	type Transformer struct {
		ID       int    `json:"id"`
		FullName string `json:"fullName"`
		Type     string `json:"type"`
	}

	trans := Transformer{FullName: "tlo1", Type: "type1"}

	row := db.QueryRow(context.Background(),
		"insert into transformer (fullName, type) values($1, $2) RETURNING id ",
		trans.FullName, trans.Type)

	var id int
	log.Error(row.Scan(&id))
	fmt.Println(id)

	err = db.Ping(context.Background())
	log.Error(err)
}

func main() {
	var msg string
	configName := flag.String("config", "./config.yaml", "Set config file path")
	_ = flag.String("log", "", "Set log path")
	flag.StringVar(&msg, "message", "hello!", "just message")
	flag.Parse()
	fmt.Printf("Begin! %v\n", configName)
	fmt.Printf("%s end\n", msg)

	hlp := helper.NewHelper("", logrus.ErrorLevel)
	conf := config.Config{}
	conf.NewConfig(hlp.Log, "./config.yaml")
	fmt.Println(conf)
	//  config.Config("./config.yaml")
	getFiles(hlp.Log, "./examples")
	// xmlparse()
	db(hlp.Log)
}
