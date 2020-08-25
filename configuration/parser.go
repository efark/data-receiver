/*
This file contains the Parser interface and three implementations.
*/

package configuration

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Parser is the interface to process input data and store it in the internal data structures.
type Parser interface {
	Parse(Configuration) error
}

// CreateParser returns an adequate parser to process the input.
func CreateParser(f, record string) Parser {
	ext := filepath.Ext(f)

	var parser Parser
	switch ext {
	case ".json":
		parser = NewJsonFileParser(f)
	case ".yaml":
		parser = NewYamlFileParser(f)
	default:
		if record != "" {
			parser = NewJsonRecordParser(record)
		} else {
			return parser
		}
	}

	return parser
}

// JsonFileParser reads a Json file and returns the configuration populated.
type JsonFileParser struct {
	filepath string
}

// NewJsonFileParser returns the JsonFileParser.
func NewJsonFileParser(f string) *JsonFileParser {
	return &JsonFileParser{filepath: f}
}

// Parse reads the file and loads the data into the Configuration interface.
func (c *JsonFileParser) Parse(conf Configuration) error {
	b, err := ioutil.ReadFile(c.filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &conf)
	if err != nil {
		return err
	}

	return nil
}

// JsonRecordParser receives only one record from CLI and uses it to populate the configuration.
type JsonRecordParser struct {
	content string
}

// NewJsonRecordParser stores a string with the Json Record.
func NewJsonRecordParser(s string) *JsonRecordParser {
	return &JsonRecordParser{content: s}
}

// Parse read the Json and loads the parameters into a struct that implements the Configuration interface.
func (c *JsonRecordParser) Parse(conf Configuration) error {
	err := json.Unmarshal([]byte(c.content), &conf)
	if err != nil {
		return err
	}

	return nil
}

// YamlFileParser reads a yaml file.
type YamlFileParser struct {
	filepath string
}

// NewYamlFileParser stores the file path for the file with the configuration.
func NewYamlFileParser(f string) *YamlFileParser {
	return &YamlFileParser{filepath: f}
}

// Parse reads the file and loads the parameters into a struct that implements the Configuration interface.
func (c *YamlFileParser) Parse(conf Configuration) error {
	b, err := ioutil.ReadFile(c.filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, conf)
	if err != nil {
		return err
	}

	return nil
}
