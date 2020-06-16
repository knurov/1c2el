package parser

import (
	"encoding/xml"
	"fmt"

	"knurov.ru/el/1c2el/helper"
)

//XMLParse parse specific xml
func XMLParse(hlp *helper.Helper, rawXML []byte) (err error) {

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
	result := TR{}
	err = xml.Unmarshal(rawXML, &result)
	if err != nil {
		return fmt.Errorf("On read xml %q", err)
	}

	for _, item := range result.Description {
		hlp.Log.Trace("Begin processing transformer %v", item.Params.Name)
		rule := findRule(hlp, item.Params.Name)
		if rule != nil {
			parseByRule(hlp, item.Params.Name, rule)
		} else {
			hlp.Log.Trace("No rule found for transformer %v", item.Params.Name)
		}
	}
	return err
}
