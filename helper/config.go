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

// FieldRule rule of parsing field
type FieldRule struct {
	Name        string `yaml:"name"`
	Rule        uint8  `yaml:"rule"`
	Transformer uint8  `yaml:"transformer"`
	Coil        uint8  `yaml:"coil"`
	Tap         uint8  `yaml:"tap"`
	Position    uint8  `yaml:"position"`
	Value       string `yaml:"value"`
	RuleRegexp  `yaml:",inline"`
}

//GroupRange describe group range
type GroupRange string

//GetRange of groups
func (groupRange *GroupRange) GetRange() (start uint8, end uint8) {
	str := fmt.Sprintf("%v", *groupRange)
	fmt.Sscanf(str, "%d..%d", &start, &end)
	return start, end
}

//End of groups
func (groupRange *GroupRange) End() (end uint8) {
	return end
}

//TransformerRule describe transformer
type TransformerRule struct {
	RuleRegexp `yaml:",inline"`
	Fields     []FieldRule `yaml:"fields"`
}

//DetailRule describe coil or tap
type DetailRule struct {
	RuleRegexp `yaml:",inline"`
	Fields     []FieldRule `yaml:"fields"`
}

//Rule desribe of rule
type Rule struct {
	Name        string `yaml:"name"`
	RuleRegexp  `yaml:",inline"`
	Transformer TransformerRule `yaml:"transformer"`
	Coils       struct {
		ParentPosition uint8 `yaml:"parentPosition"`
		StartPosition  uint8 `yaml:"startPosition"`
		EndPosition    uint8 `yaml:"endPosition"`
		RuleRegexp     `yaml:",inline"`
		Fields         []FieldRule `yaml:"fields"`
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
	Rules []Rule `yaml:"правило"`

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

	for ruleIndex, item := range config.Rules {
		err = config.Rules[ruleIndex].RegexpCompile()
		config.Loger.Fatal("On Compile Transformer regexp %v", err)

		for coilIndex := range item.Coils.Fields {
			err = config.Rules[ruleIndex].Coils.Fields[coilIndex].RegexpCompile()
			config.Loger.Fatal("On Compile Coils regexp %v", err)
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
}

//Destroy Destroyng config (closing opened files, databases and etc)
func (config *Config) Destroy() {
	config.Database.Pool.Close()
}
