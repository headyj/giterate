package command

import (
	"encoding/json"
	"giterate/pkg/entities"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/go-yaml/yaml"
)

func initConf(arguments *entities.Arguments) []entities.Service {
	var file []byte
	var err error
	var configFilePath string
	if arguments.ConfigFile != "" {
		_, configFile := os.Stat(arguments.ConfigFile)
		if os.IsNotExist(configFile) {
			log.Fatalf("Configuration file not found\n")
		}
		configFilePath = arguments.ConfigFile
	} else {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Cannot get current user:%s\n", err)
		}
		_, json := os.Stat(usr.HomeDir + "/" + ".giterate/config.json")
		_, yaml := os.Stat(usr.HomeDir + "/" + ".giterate/config.yaml")
		if os.IsNotExist(json) {
			if os.IsNotExist(yaml) {
				log.Fatalf("Configuration file not found\n")
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
		log.Fatalf("Configuration file cannot be parsed:%s\n", err)
	}
	switch arguments.LogLevel {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	return services
}
