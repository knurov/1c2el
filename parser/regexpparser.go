package parser

import (
	"knurov.ru/el/1c2el/db"
	"knurov.ru/el/1c2el/helper"
)

func parseByRule(hlp *helper.Helper, fullName string, rule *helper.Rule) {
	hlp.Log.Trace("Use rule %v", rule.Name)

	// Поля трансформатора #####################
	// Определение источников для полей трансформатора #############
	transformer := make(map[string]string)
	ruleResult := rule.FindAllGroup(fullName)
	var transformerResult []string
	if rule.Transformer.RegexpCompiled != nil {
		transformerResult = rule.Transformer.FindAllGroup(fullName)
	}
	// ############# Определение источников для полей трансформатора

	// Цикл сбора полей трансформатора #############
	for _, field := range rule.Transformer.Fields {
		name, value, err := field.GetFieldMap(ruleResult, transformerResult, nil, nil)
		hlp.Log.Error(err)
		transformer[name] = value
	}
	// ############# Цикл сбора полей трансформатора
	// ##################### Поля трансформатора

	// Поля отпаек #####################
	//  Цикл обмотк ################
	coils := rule.Coils.Rule.GetRange(ruleResult)

	// Rule        GroupRange `yaml:"rule"`
	// Transformer GroupRange `yaml:"transformer"`
	// Position    GroupRange `yaml:"position"`

	// Separator   string     `yaml:"separator"`
	// RuleRegexp  `yaml:",inline"`
	// Taps        []TapRule `yaml:"taps"`

	for index, value := range coils {
		taps := rule.Coils.Taps.coil.GetRange(value)

		//  Цикл отпаек ################

		// ################# Цикл отпаек
		// ################# Цикл обмотк
	}
	// ##################### Поля отпаек

	// вставка данных в BD #################
	transformeID := db.Transformer(hlp, transformer)
	coil := make(map[string]string)
	coil["transformeID"] = string(transformeID)
	// ################# вставка данных в BD

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
