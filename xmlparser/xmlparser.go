package xmlparser

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"knurov.ru/el/1c2el/db"
	"knurov.ru/el/1c2el/helper"
)

// func parseTLO10(fullName string)
// https://yourbasic.org/golang/regexp-cheat-sheet/
func parseFullName(hlp *helper.Helper, fullName string) {
	isTlo10, err := regexp.Compile("^ТЛО-10")
	// 'ТЛО-10_М1ACE-0.2SFS7/0.5FS7/10P10-10/10/40-150(300)-150(300)-300/5 У2 б 31.5кА'
	// tlo10, err := regexp.Compile("^(ТЛО-10)_(.+)-(.+)/(.+)/(.+)/(.+)/(.+)/(.+) (.+) (.+) (.+)")
	tlo10, err := regexp.Compile(`(?P<short>.+?)_(?P<prop>.+?)-`)
	hlp.Log.Fatal(err)
	if isTlo10.MatchString(fullName) {
		result := tlo10.FindStringSubmatch(fullName)
		fmt.Println(tlo10.SubexpNames())
		// for _, item := range result {
		// 	fmt.Println(item)
		// }
		// fmt.Println(result)
		if len(result) == 3 {
			db.Transformer(hlp, result[1], result[2])
		} else {
			hlp.Log.Trace("No result")
		}
	}
}

//XMLParse parse specific xml
func XMLParse(hlp *helper.Helper, fileName string) {

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
	hlp.Log.Fatal(err)
	result := TR{}
	xmlData, err := ioutil.ReadAll(file)
	hlp.Log.Fatal(err)
	err = xml.Unmarshal(xmlData, &result)

	for _, item := range result.Description {
		// fmt.Printf("%v - Serial number %v\n", item.Params.Name, item.Number)
		parseFullName(hlp, item.Params.Name)
	}
}
