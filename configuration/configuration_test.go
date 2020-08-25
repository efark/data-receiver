package configuration_test

import (
	"errors"
	"fmt"
	"github.com/efark/data-receiver/configuration"
	"os"
	"path"
	"strings"
	"testing"
)

var baseFilepath = "./test_config"
var jsonContent = `{"services": {"test_service": {"extractor": {"type": "HeaderExtractor"}, "authenticator": {"type": "Signer", "parameters": {"Key": "magicKey"}}, "writer": {"type": "MemoryWriter"}}}}`
var yamlContent = `services: 
  test_service:
    extractor: 
      type: HeaderExtractor
    authenticator: 
      type: Signer
      parameters:
        Key: "magicKey"
    writer: 
      type: MemoryWriter`

func TestNewServiceMap(t *testing.T) {
	var conf configuration.Configuration
	tmap := configuration.NewServiceMap()
	conf = tmap

	testConf := configuration.NewServiceConfig(
		&configuration.SimpleConfig{Class: "HeaderExtractor", Parameters: make(map[string]string)},
		&configuration.SimpleConfig{Class: "Signer", Parameters: map[string]string{"Key": "magicKey"}},
		&configuration.SimpleConfig{Class: "MemoryWriter", Parameters: make(map[string]string)})

	err := conf.Add("test_service", testConf)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if len(conf.List()) != 1 {
		t.Error(errors.New(fmt.Sprintf("Expected len(conf.List()): - %d Received: %d", 1, len(conf.List()))))
		t.FailNow()
	}

	serv, err := conf.Get("test_service")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if serv.ExtConfig.Class != "HeaderExtractor" {
		t.Error(errors.New(fmt.Sprintf("Expected serv.AuthConfig.Class: - %q Received: %q", "HeaderExtractor", serv.AuthConfig.Class)))
		t.FailNow()
	}

	if serv.AuthConfig.Class != "Signer" {
		t.Error(errors.New(fmt.Sprintf("Expected serv.AuthConfig.Class: - %q Received: %q", "Signer", serv.AuthConfig.Class)))
		t.FailNow()
	}

	if serv.AuthConfig.Parameters["Key"] != "magicKey" {
		t.Error(errors.New(fmt.Sprintf("Expected serv.AuthConfig.Parameters[\"Key\"]: - %q Received: %q", "magicKey", serv.AuthConfig.Parameters["Key"])))
		t.FailNow()
	}

	if serv.AuthConfig.Class != "Signer" {
		t.Error(errors.New(fmt.Sprintf("Expected serv.AuthConfig.Class: - %q Received: %q", "Signer", serv.AuthConfig.Class)))
		t.FailNow()
	}
}

func TestJsonRecordParser(t *testing.T) {

	var conf configuration.Configuration
	realConf := configuration.NewServiceMap()
	conf = realConf

	parser := configuration.CreateParser("", jsonContent)
	err := parser.Parse(conf)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	list := conf.List()
	if len(list) != 1 {
		t.Error(fmt.Sprintf("Error: Expected len(list) == %d, obtained %d.", 1, len(list)))
		t.FailNow()
	}

	serv, err := conf.Get(list[0])
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if list[0] != "test_service" {
		t.Error(fmt.Sprintf("Error: Expected serv.Name == %s, obtained %s.", "test_service", list[0]))
		t.FailNow()
	}

	if serv.ExtConfig.Class != "HeaderExtractor" {
		t.Error(fmt.Sprintf("Error: Expected serv.ExtConfig.Class() == %s, obtained %s.", "HeaderExtractor", serv.ExtConfig.Class))
		t.FailNow()
	}

	if serv.AuthConfig.Class != "Signer" {
		t.Error(fmt.Sprintf("Error: Expected serv.AuthConfig.Class() == %s, obtained %s.", "Signer", serv.AuthConfig.Class))
		t.FailNow()
	}

	if serv.AuthConfig.Parameters["Key"] != "magicKey" {
		t.Error(fmt.Sprintf("Error: Expected serv.AuthConfig.Parameters['Key'] == %s, obtained %s.", "magicKey", serv.AuthConfig.Parameters["Key"]))
		t.FailNow()
	}

	if serv.WriConfig.Class != "MemoryWriter" {
		t.Error(fmt.Sprintf("Error: Expected serv.WriConfig.Class() == %s, obtained %s.", "MemoryWriter", serv.WriConfig.Class))
		t.FailNow()
	}
}

func setupTest(t *testing.T, filepath string) func() {
	if fileExists(filepath) {
		err := os.Remove(filepath)
		if err != nil {
			t.Error(err.Error())
		}
	}

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	ext := path.Ext(filepath)
	var content string
	if ext == ".json" {
		content = jsonContent
	} else {
		content = yamlContent
	}

	_, err = file.WriteString(content)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	err = file.Close()
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	return func() {
		if fileExists(filepath) {
			err := os.Remove(filepath)
			if err != nil {
				t.Error(err.Error())
			}
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func TestJsonFileParser(t *testing.T) {
	f := strings.Join([]string{baseFilepath, "json"}, ".")

	defer setupTest(t, f)()

	var conf configuration.Configuration

	realConf := configuration.NewServiceMap()
	conf = realConf

	parser := configuration.CreateParser(f, "")
	err := parser.Parse(conf)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	list := conf.List()
	if len(list) != 1 {
		t.Error(fmt.Sprintf("Error: Expected len(list) == %d, obtained %d.", 1, len(list)))
		t.FailNow()
	}

	serv, err := conf.Get(list[0])
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if list[0] != "test_service" {
		t.Error(fmt.Sprintf("Error: Expected serv.Name == %s, obtained %s.", "test_service", list[0]))
		t.FailNow()
	}

	if serv.ExtConfig.Class != "HeaderExtractor" {
		t.Error(fmt.Sprintf("Error: Expected serv.ExtConfig.Class() == %s, obtained %s.", "HeaderExtractor", serv.ExtConfig.Class))
		t.FailNow()
	}

	if serv.AuthConfig.Class != "Signer" {
		t.Error(fmt.Sprintf("Error: Expected serv.AuthConfig.Class() == %s, obtained %s.", "Signer", serv.AuthConfig.Class))
		t.FailNow()
	}

	if serv.AuthConfig.Parameters["Key"] != "magicKey" {
		t.Error(fmt.Sprintf("Error: Expected serv.AuthConfig.Parameters['Key'] == %s, obtained %s.", "magicKey", serv.AuthConfig.Parameters["Key"]))
		t.FailNow()
	}

	if serv.WriConfig.Class != "MemoryWriter" {
		t.Error(fmt.Sprintf("Error: Expected serv.WriConfig.Class() == %s, obtained %s.", "MemoryWriter", serv.WriConfig.Class))
		t.FailNow()
	}
}

func TestYamlFileParser(t *testing.T) {
	f := strings.Join([]string{baseFilepath, "yaml"}, ".")
	defer setupTest(t, f)()

	var conf configuration.Configuration
	realConf := configuration.NewServiceMap()
	conf = realConf

	parser := configuration.CreateParser(f, "")
	err := parser.Parse(conf)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	list := conf.List()
	if len(list) != 1 {
		t.Error(fmt.Sprintf("Error: Expected len(list) == %d, obtained %d.", 1, len(list)))
		t.FailNow()
	}

	serv, err := conf.Get(list[0])
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	// fmt.Printf("%+v\n", serv)

	if list[0] != "test_service" {
		t.Error(fmt.Sprintf("Error: Expected serv.Name == %s, obtained %s.", "test_service", list[0]))
		t.FailNow()
	}

	if serv.ExtConfig.Class != "HeaderExtractor" {
		t.Error(fmt.Sprintf("Error: Expected serv.ExtConfig.Class() == %s, obtained %s.", "HeaderExtractor", serv.ExtConfig.Class))
		t.FailNow()
	}

	if serv.AuthConfig.Class != "Signer" {
		t.Error(fmt.Sprintf("Error: Expected serv.AuthConfig.Class() == %s, obtained %s.", "Signer", serv.AuthConfig.Class))
		t.FailNow()
	}

	if serv.AuthConfig.Parameters["Key"] != "magicKey" {
		t.Error(fmt.Sprintf("Error: Expected serv.AuthConfig.Parameters['Key'] == %s, obtained %s.", "magicKey", serv.AuthConfig.Parameters["Key"]))
		t.FailNow()
	}

	if serv.WriConfig.Class != "MemoryWriter" {
		t.Error(fmt.Sprintf("Error: Expected serv.WriConfig.Class() == %s, obtained %s.", "MemoryWriter", serv.WriConfig.Class))
		t.FailNow()
	}
}
