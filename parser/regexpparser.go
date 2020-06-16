package parser

import (
	"knurov.ru/el/1c2el/db"
	"knurov.ru/el/1c2el/helper"
)

func parseByRule(hlp *helper.Helper, fullName string, rule *helper.Rule) {
	hlp.Log.Trace("Use rule %v", rule.Name)
	transformer := make(map[string]string)
	ruleResult := rule.RegexpCompiled.FindStringSubmatch(fullName)
	var transformerResult []string
	if rule.Transformer.RegexpCompiled != nil {
		transformerResult = rule.Transformer.RegexpCompiled.FindStringSubmatch(fullName)
	}

	for _, field := range rule.Transformer.Fields {
		if field.Rule != 0 {
			transformer[field.Name] = ruleResult[field.Rule]
		} else if field.Transformer != 0 {
			if len(transformerResult) < int(field.Transformer) {
				hlp.Log.Error("Count of groups '%v' smallest than index of field '%v'", len(transformerResult), field.Transformer)
			} else {
				transformer[field.Name] = transformerResult[field.Transformer]
			}
		}
	}

	transformeID := db.Transformer(hlp, transformer)
	coil := make(map[string]string)
	coil["transformeID"] = string(transformeID)
	// for _, field := range rule.Coils.Fields {
	// 	transformer[field.Name] = ruleResult[field.Position]
	// }

}

func findRule(hlp *helper.Helper, fullName string) *helper.Rule {
	// https://yourbasic.org/golang/regexp-cheat-sheet/
	// 'ТЛО-10_М1ACE-0.2SFS7/0.5FS7/10P10-10/10/40-150(300)-150(300)-300/5 У2 б 31.5кА'
	// tlo10, err := regexp.Compile("^(ТЛО-10)_(.+)-(.+)/(.+)/(.+)/(.+)/(.+)/(.+) (.+) (.+) (.+)")

	// fmt.Println(hlp.Conf.Rules)
	for _, rule := range hlp.Conf.Rules {
		if rule.RegexpCompiled != nil && rule.RegexpCompiled.MatchString(fullName) {
			hlp.Log.Trace("Found rule %v for %v", rule.Name, fullName)
			return &rule
		}
		hlp.Log.Trace("Skip rule %v for %v", rule.Name, fullName)
	}
	return nil
}
