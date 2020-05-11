package xmlparser

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"knurov.ru/el/1c2el/db"
	"knurov.ru/el/1c2el/helper"
)

func parseByRule(hlp *helper.Helper, fullName string) {
	// https://yourbasic.org/golang/regexp-cheat-sheet/
	// 'ТЛО-10_М1ACE-0.2SFS7/0.5FS7/10P10-10/10/40-150(300)-150(300)-300/5 У2 б 31.5кА'
	// tlo10, err := regexp.Compile("^(ТЛО-10)_(.+)-(.+)/(.+)/(.+)/(.+)/(.+)/(.+) (.+) (.+) (.+)")

	for _, rule := range hlp.Conf.Rules {

		if rule.RegexpCompiled != nil && rule.RegexpCompiled.MatchString(fullName) {
			result := rule.RegexpCompiled.FindStringSubmatch(fullName)

			transformer := make(map[string]string)
			for _, transformerField := range rule.Transformer {
				transformer[transformerField.Field] = result[transformerField.Position]
			}

			db.Transformer(hlp, transformer)
		} else {
			hlp.Log.Trace("Skip rule %v for %v", rule.Name, fullName)
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
		hlp.Log.Trace("Process transformer %v", item.Params.Name)
		parseByRule(hlp, item.Params.Name)
	}
}
