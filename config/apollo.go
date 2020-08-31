package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/apollo.v0"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func LoadFromApollo() {
	LoadFromEnvFile()
	if IsApolloEnable() {

		if err := apollo.StartWithConf(&config.Apollo.Settings); err != nil {
			fmt.Printf("load apollo config error: %s\n", err)
			os.Exit(1)
			return
		}
		readConfig(config)

		// hot loading refresh config
		go func() {
			for {
				event := apollo.WatchUpdate()
				changeEvent := <-event
				bytes, _ := json.Marshal(changeEvent)
				fmt.Println("event:", string(bytes))
				readConfig(config)
				readConfig(customConfig)

			}
		}()
	}
}

var requireLoadCustomConfig bool
var customConfig interface{}

//LoadCustomFromFile Load custom config from apollo, save to custom config
func LoadCustomFromApollo(customCfg interface{}) error {
	if !IsApolloEnable() {

		return errors.New("apollo is not enabled")
	}
	requireLoadCustomConfig = true
	readConfig(customCfg)
	customConfig = customCfg
	return nil
}
func readConfig(cfg interface{}) {
	for _, namespace := range config.Apollo.Settings.Namespaces {
		ymlTxt := apollo.GetNameSpaceContent(namespace, "")
		if err := yaml.NewDecoder(strings.NewReader(ymlTxt)).Decode(cfg); err != nil {
			fmt.Printf("load apollo config error: %s\n", err)
		}
	}
}
