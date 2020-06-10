package parser

import (
	"knurov.ru/el/1c2el/db"
	"knurov.ru/el/1c2el/helper"
)

func parseByRule(hlp *helper.Helper, fullName string, rule *helper.Rule) {
	hlp.Log.Trace("Use rule %v", rule.Name)
	result := rule.RegexpCompiled.FindStringSubmatch(fullName)

	transformer := make(map[string]string)
	for _, field := range rule.Transformer.Fields {
		transformer[field.Name] = result[field.Position]
	}

	transformeID := db.Transformer(hlp, transformer)
	coil := make(map[string]string)
	coil["transformeID"] = string(transformeID)
	for _, field := range rule.Coils.Fields {
		transformer[field.Name] = result[field.Position]
	}

}
