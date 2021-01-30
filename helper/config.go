package helper

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"

	yamlconvert "github.com/ghodss/yaml"
	"github.com/jackc/pgx/v4/pgxpool"
	logrus "github.com/sirupsen/logrus"
	gojsonschema "github.com/xeipuuv/gojsonschema"

	"gopkg.in/yaml.v2"
)

//RuleRegexp base rule
type RuleRegexp struct {
	Regexp         string         `yaml:"regexp"`
	RegexpCompiled *regexp.Regexp `yaml:"-"`
}

//RegexpCompile compile
func (ruleRegexp *RuleRegexp) RegexpCompile() (err error) {
	ruleRegexp.RegexpCompiled, err = regexp.Compile(ruleRegexp.Regexp)
	return err
}

//FindAllGroup retur group array
func (ruleRegexp *RuleRegexp) FindAllGroup(value string) (result []string) {
	submatch := ruleRegexp.RegexpCompiled.FindAllStringSubmatch(value, -1)
	for _, match := range submatch {
		for _, group := range match[1:] {
			if group != "" {
				result = append(result, group)
			}
		}
	}
	return result
}

// FieldRule rule of parsing field
type FieldRule struct {
	Rule        uint8  `yaml:"rule"`
	Transformer uint8  `yaml:"transformer"`
	Coil        uint8  `yaml:"coil"`
	Tap         uint8  `yaml:"tap"`
	Name        string `yaml:"name"`
	Position    uint8  `yaml:"position"`
	Value       string `yaml:"value"`
	RuleRegexp  `yaml:",inline"`
}

//GetFieldMap get parsed result
func (field *FieldRule) GetFieldMap(rule []string, transformer []string, coil []string, tap []string) (name string, value string, err error) {

	// var allGroup string

	if field.Rule > 0 {
		return field.Name, rule[field.Rule-1], nil
	} else if field.Transformer > 0 {
		if len(transformer) < int(field.Transformer) {
			return field.Name, "", fmt.Errorf("Count of groups '%v' smallest than index of field '%v'", len(transformer), field.Transformer)
		} else {
			return field.Name, transformer[field.Transformer-1], nil
		}
	} else if field.Coil > 0 {
		if len(coil) < int(field.Coil) {
			return field.Name, "", fmt.Errorf("Count of groups '%v' smallest than index of field '%v'", len(coil), field.Coil)
		} else {
			return field.Name, coil[field.Coil-1], nil
		}
	} else if field.Tap > 0 {
		if len(coil) < int(field.Coil) {
			return field.Name, "", fmt.Errorf("Count of groups '%v' smallest than index of field '%v'", len(coil), field.Coil)
		} else {
			return field.Name, coil[field.Coil-1], nil
		}
	} else if field.Value != "" {
		return field.Name, field.Value, nil
	}
	return "field.Value", "", fmt.Errorf("No match value by rule")

}

//GroupRange describe group range
type GroupRange string

//GetIndexGap of groups
func (groupRange *GroupRange) GetIndexGap() (start int8, end int8) {
	fmt.Sscanf(string(*groupRange), "%d..%d", &start, &end)
	return start, end
}

//GetRange of groups
func (groupRange *GroupRange) GetRange(groups []string) (items []string) {
	from, to := groupRange.GetIndexGap()
	for index, value := range groups {
		if index > int(to) && to > -1 || index > len(groups)-int(to) {
			break
		}
		if index > int(from) {
			items[len(items)] = value
		}
	}
	return items
}

//TransformerRule describe transformer
type TransformerRule struct {
	RuleRegexp `yaml:",inline"`
	Fields     []FieldRule `yaml:"fields"`
}

// //DetailRule describe coil or tap
// type DetailRule struct {
// 	RuleRegexp `yaml:",inline"`
// 	Fields     []FieldRule `yaml:"fields"`
// }

//TapRule Tap Rule
type TapRule struct {
	Rule        GroupRange `yaml:"rule"`
	Transformer GroupRange `yaml:"transformer"`
	Coil        GroupRange `yaml:"coil"`
	Position    GroupRange `yaml:"position"`
	Separator   string     `yaml:"separator"`
	RuleRegexp  `yaml:",inline"`
	Fields      []FieldRule `yaml:"fields"`
}

//Rule desribe of rule
type Rule struct {
	Name        string `yaml:"name"`
	RuleRegexp  `yaml:",inline"`
	Transformer TransformerRule `yaml:"transformer"`
	Coils       struct {
		Rule        GroupRange `yaml:"rule"`
		Transformer GroupRange `yaml:"transformer"`
		Position    GroupRange `yaml:"position"`
		Separator   string     `yaml:"separator"`
		RuleRegexp  `yaml:",inline"`
		Taps        []TapRule `yaml:"taps"`
	} `yaml:"coils"`
}

// Config - Содержит параметры приложения
type Config struct {
	Database struct {
		DryRun bool   `yaml:"dryRun"`
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Db     string `yaml:"db"`
		User   string `yaml:"user"`
		Passwd string `yaml:"passwd"`
		Pool   *pgxpool.Pool
	} `yaml:"database"`
	Files struct {
		Src          string `yaml:"src"`
		Done         string `yaml:"done"`
		Error        string `yaml:"error"`
		NameTemplate string `yaml:"nameTemplate"`
		NameRegexp   *regexp.Regexp
	}
	Log struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"log"`

	//PaseRules rule parsing rule
	Rules []Rule `yaml:"rules"`

	Loger *Loger
}

//NewConfig -  КОнструктор
func NewConfig(fileName string, dryRun bool) (config *Config) {
	config = new(Config)
	config.init(NewLoger(), fileName, dryRun)
	return config
}

func (config *Config) validate(configFile []byte) {
	// https: //json-schema.org/learn/miscellaneous-examples.html
	schema := gojsonschema.NewStringLoader(`
		{
			"required": [ "database", "files" ],
			"properties": {
				"database": {
					"type": "object",
					"required": [ "host", "port" ],
					"dryRun": {
						"type": "boolean"
					},	
					"host": {
						"type": "string"
					},
					"port": {
						"type": "int"
					}
				},
				"files": {
					"type": "object",
					"src": {
						"type": "string"
					} 
				}
			}
		} 
	`)

	configJSON, err := yamlconvert.YAMLToJSON(configFile)
	config.Loger.Fatal("On config validadion %v", err)
	validationResult, err := gojsonschema.Validate(schema, gojsonschema.NewBytesLoader(configJSON))
	config.Loger.Fatal("Validadion error: %v", err)

	if !validationResult.Valid() {
		config.Loger.Fatal("Validadion error: %v", validationResult.Errors())
	}

}

func (config *Config) setDefaults() {
	config.Database.Port = 5432
	config.Database.Host = "localhost"
	config.Log.Level = logrus.ErrorLevel.String()
	config.Files.NameTemplate = ".*\\.xml"
}

func (config *Config) setValuesRules() {
	var err error

	for ruleIndex := range config.Rules {
		rule := &config.Rules[ruleIndex]
		err = rule.RegexpCompile()
		config.Loger.Fatal("On Compile Rule regexp %v", err)

		err = rule.Transformer.RegexpCompile()
		config.Loger.Fatal("On Compile Transformer regexp %v", err)

		for transformerFieldIndex := range rule.Transformer.Fields {
			err = rule.Transformer.Fields[transformerFieldIndex].RegexpCompile()
			config.Loger.Fatal("On Compile Transformer fields regexp %v", err)
		}

		for tapIndex := range rule.Coils.Taps {
			tap := &rule.Coils.Taps[tapIndex]
			err = tap.RegexpCompile()
			config.Loger.Fatal("On Compile Coils regexp %v", err)
			for tapFieldIndex := range tap.Fields {
				err = tap.Fields[tapFieldIndex].RegexpCompile()
				config.Loger.Fatal("On Compile Coils regexp %v", err)
			}

		}

	}

}

func (config *Config) setValues(dryRun bool) {
	var err error
	config.Database.DryRun = config.Database.DryRun || dryRun
	config.Files.NameRegexp, err = regexp.Compile(config.Files.NameTemplate)
	config.Loger.Fatal("On Compile FIleName regexp %v", err)

	level, err := logrus.ParseLevel(config.Log.Level)
	config.Loger.Fatal("On parse Log Level %v", err)
	config.Loger.Init(config.Log.File, level)

	// postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]
	connectionString := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v",
		config.Database.User,
		config.Database.Passwd,
		config.Database.Host,
		config.Database.Port,
		config.Database.Db)
	dbConfig, err := pgxpool.ParseConfig(connectionString)
	config.Loger.Fatal("On create config url - %v", err)
	config.Database.Pool, err = pgxpool.ConnectConfig(context.Background(), dbConfig)
	config.Loger.Fatal("On connect to db - %v", err)

	config.setValuesRules()
}

//init -  КОнструктор
func (config *Config) init(log *Loger, fileName string, dryRun bool) {
	config.Loger = log
	configFile, err := ioutil.ReadFile(fileName)
	config.Loger.Fatal("On read config file: %v", err)
	config.validate(configFile)
	config.setDefaults()
	err = yaml.Unmarshal(configFile, &config)
	config.Loger.Fatal("On parse "+fileName+" v%", err)
	config.setValues(dryRun)
	if config.Database.DryRun {
		fmt.Printf("%+v\n", config)
	}
}

//Destroy Destroyng config (closing opened files, databases and etc)
func (config *Config) Destroy() {
	config.Database.Pool.Close()
}
