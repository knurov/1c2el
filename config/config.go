package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	yamlconvert "github.com/ghodss/yaml"
	gojsonschema "github.com/xeipuuv/gojsonschema"
	"knurov.ru/el/1c2el/helper"

	"gopkg.in/yaml.v2"
)

// Config - Содержит параметры приложения
type Config struct {
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"database"`
	Log struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"log"`
}

//NewConfig -  КОнструктор
func (config *Config) NewConfig(fileName string) {
	// https://json-schema.org/learn/miscellaneous-examples.html
	schema := gojsonschema.NewStringLoader(`
		{
			"required": [ "database" ],
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
				}
			}
		} 
	`)

	configFile, err := ioutil.ReadFile(fileName)
	helper.LogFatal(err)
	configJSON, err := yamlconvert.YAMLToJSON(configFile)
	helper.LogFatal(err)
	validationResult, err := gojsonschema.Validate(schema, gojsonschema.NewBytesLoader(configJSON))
	helper.LogFatal(err)

	if !validationResult.Valid() {
		log.Fatal(validationResult.Errors())
	}

	err = yaml.Unmarshal(configFile, &config)
	helper.LogFatal(err)

	config.Database.Port = 5432
	config.Database.Host = "localhost"
	config.Log.Level = log.ErrorLevel.String()
}
