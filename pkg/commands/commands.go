package command

import (
	"encoding/json"
	"giterate/pkg/entities"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

func initConf(arguments *Arguments) []entities.Service {
	var file []byte
	var err error
	var configFilePath string
	if arguments.ConfigFile != "" {
		_, configFile := os.Stat(arguments.ConfigFile)
		if os.IsNotExist(configFile) {
			log.Fatalf("Configuration file not found")
		}
		configFilePath = arguments.ConfigFile
	} else {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Cannot get current user:%s", err)
		}
		_, json := os.Stat(usr.HomeDir + "/" + ".giterate/config.json")
		_, yaml := os.Stat(usr.HomeDir + "/" + ".giterate/config.yaml")
		if os.IsNotExist(json) {
			if os.IsNotExist(yaml) {
				log.Fatalf("Configuration file not found")
			}
			configFilePath = usr.HomeDir + "/" + ".giterate/config.yaml"
		} else {
			configFilePath = usr.HomeDir + "/" + ".giterate/config.json"
		}
	}
	file, _ = ioutil.ReadFile(configFilePath)

	var services []entities.Service
	if filepath.Ext(configFilePath) == ".yaml" {
		err = yaml.Unmarshal(file, &services)
	} else {
		err = json.Unmarshal(file, &services)
	}
	if err != nil {
		log.Fatalf("File cannot be parsed:%s", err)
	}
	return services
}
