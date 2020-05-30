package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"knurov.ru/el/1c2el/helper"
	"knurov.ru/el/1c2el/xmlparser"
)

func readFile(hlp *helper.Helper, fileName string) (err error) {
	hlp.Log.Info("Processing file %v", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("On open file %q", err)
	}
	defer file.Close()
	hlp.Log.Fatal(err)
	rawXML, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("On read file %q", err)
	}
	err = xmlparser.XMLParse(hlp, rawXML)
	return err
}

func moveTo(logger *helper.Loger, fullFileName string, dst string) {
	if dst == "" {
		return
	}
	dstFullName := filepath.Join(dst, filepath.Base(fullFileName))
	logger.Trace("Move %q to %q", fullFileName, dst)
	err := os.Rename(fullFileName, dstFullName)
	logger.Error(err)
}

func readDir(hlp *helper.Helper) {
	hlp.Log.Info("Processing dir %v", hlp.Conf.Files.Src)
	files, err := ioutil.ReadDir(hlp.Conf.Files.Src)
	hlp.Log.Fatal("On read the dir %v", err)
	for _, item := range files {
		if item.IsDir() {
			continue
		}
		fullFileNaem := filepath.Join(hlp.Conf.Files.Src, item.Name())
		err = readFile(hlp, fullFileNaem)
		hlp.Log.Error(err)
		if err == nil {
			moveTo(hlp.Log, fullFileNaem, hlp.Conf.Files.Done)
		} else {
			moveTo(hlp.Log, fullFileNaem, hlp.Conf.Files.Error)
		}
	}
	hlp.Log.Info("Done read dir %v", hlp.Conf.Files.Src)
}
