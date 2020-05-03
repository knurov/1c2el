package config

import (
	"io/ioutil"
	"regexp"

	"knurov.ru/el/1c2el/helper"

	yamlconvert "github.com/ghodss/yaml"
	gojsonschema "github.com/xeipuuv/gojsonschema"

	"gopkg.in/yaml.v2"
)

// Config - Содержит параметры приложения
type Config struct {
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
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
	Loger *helper.Loger
}

//NewConfig -  КОнструктор
func NewConfig(log *helper.Loger, fileName string) (config *Config) {
	config = new(Config)
	config.init(log, fileName)
	return config
}

func (config *Config) validate(configFile []byte) {
	// https://json-schema.org/learn/miscellaneous-examples.html
	schema := gojsonschema.NewStringLoader(`
		{
			"required": [ "database", "files" ],
			"properties": {
				"database": {
					"type": "object",
					"required": [ "host", "port" ],
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
	config.Loger.Fatal("Validadion error: %v", err)
	validationResult, err := gojsonschema.Validate(schema, gojsonschema.NewBytesLoader(configJSON))
	config.Loger.Fatal("Validadion error: %v", err)

	if !validationResult.Valid() {
		config.Loger.Fatal("Validadion error: %v", validationResult.Errors())
	}

}

func (config *Config) setDefaults() {
	config.Database.Port = 5432
	config.Database.Host = "localhost"
	// config.Log.Level = log.ErrorLevel.String()
	config.Files.NameTemplate = ".*\\.xml"

}

//NewConfig -  КОнструктор
func (config *Config) init(log *helper.Loger, fileName string) {
	config.Loger = log
	configFile, err := ioutil.ReadFile(fileName)
	config.Loger.Fatal(err)
	config.validate(configFile)

	config.setDefaults()

	err = yaml.Unmarshal(configFile, &config)
	config.Loger.Fatal(err)

	config.Files.NameRegexp, err = regexp.Compile(config.Files.NameTemplate)
	config.Loger.Fatal(err)

}
