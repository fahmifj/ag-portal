package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/fahmifj/ag-portal/logger"
	"github.com/joho/godotenv"
)

var (
	config         = make(map[string]string)
	subscriptionID = "subscriptionId"
	rgName         string
)

func parseConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	// This is a workaround to get full path of AZURE_AUTH_LOCATION (assuming the config is '~/.azure')
	// Since godotenv seems can't parse $HOME in .env
	_ = os.Setenv("AZURE_AUTH_LOCATION", path.Join(os.Getenv("HOME"), os.Getenv("AZURE_AUTH_LOCATION")))

	data, err := ioutil.ReadFile(os.Getenv("AZURE_AUTH_LOCATION"))
	if err != nil {
		return err
	}

	kv := make(map[string]interface{})
	err = json.Unmarshal(data, &kv)
	if err != nil {
		return err
	}

	for k, v := range kv {
		// NOTE:
		// Currently I only need the subscription ID
		// Auth is already handled by autorest
		if k == subscriptionID {
			config[k] = v.(string)
		}
	}

	rgName = os.Getenv("RG_NAME")
	return nil
}

func Load() {
	err := parseConfig()
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}
}

func SubscriptionID() string {
	return config[subscriptionID]
}

func RG() string {
	return rgName
}
