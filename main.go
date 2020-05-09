package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"knurov.ru/el/1c2el/xmlparser"

	"knurov.ru/el/1c2el/helper"
)

func readFiles(hlp *helper.Helper) {
	files, err := ioutil.ReadDir(hlp.Conf.Files.Src)
	hlp.Log.Fatal("On read the dir %v", err)
	for _, item := range files {
		if !item.IsDir() {
			hlp.Log.Debug(fmt.Sprintf("Processing file %v", item.Name()))
			xmlparser.XMLParse(hlp, filepath.Join(hlp.Conf.Files.Src, item.Name()))
		}

	}
}

func main() {
	configName := flag.String("config", "./config.yaml", "Set config file path")
	flag.Parse()
	fmt.Printf("Starting with! %v\n", *configName)

	hlp := helper.NewHelper(*configName)
	hlp.Log.Debug("Use config - %v", *configName)
	defer hlp.Destroy()
	readFiles(hlp)
}
