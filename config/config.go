package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var POSTGRES_HOST_SLAVE string
var POSTGRES_PORT_SLAVE int
var POSTGRES_USERNAME_SLAVE string
var POSTGRES_PASSWORD_SLAVE string
var POSTGRES_DATABASE_SLAVE string
var POSTGRES_SSL_MODE_SLAVE string

var ELASTICSEARCH_CLUSTER string
var ELASTICSEARCH_USERNAME string
var ELASTICSEARCH_PASSWORD string
var ELASTICSEARCH_PRODUCT_INDEX string

// isFileExist check if the file exist on the given file path
func isFileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return false
}

// SetConfig will set config from .env file if it's exist
//
// Otherwise it will set from system's ENV variables.
// Filename should be and env file: .env or .env.* file.
// Dirpath should be in this format: /some/dirpath.
func SelectConfig(environment string) error {

	switch environment {
	case "local":
		SetConfig(".", "local.env")
		return nil
	case "development":
		SetConfig(".", "development.env")
		return nil
	case "staging":
		SetConfig(".", "staging.env")
		return nil
	default:
		return errors.New("environment is not defined")
	}

}

// SetConfig will set config from .env file if it's exist
//
// Otherwise it will set from system's ENV variables.
// Filename should be and env file: .env or .env.* file.
// Dirpath should be in this format: /some/dirpath.
func SetConfig(dirpath string, filename string) {
	filePath := filepath.Join(dirpath, filename)
	fileExist := isFileExist(filePath)

	if fileExist {
		viper.AddConfigPath(dirpath)
		viper.SetConfigFile(filePath)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("error reading config file: %+v", err)
		}
		reloadConfig()
		return
	}

	log.Fatalf("file %s not found", filePath)

}

func reloadConfig() {
	POSTGRES_HOST_SLAVE = viper.GetString("POSTGRES_HOST_SLAVE")
	POSTGRES_PORT_SLAVE = viper.GetInt("POSTGRES_PORT_SLAVE")
	POSTGRES_USERNAME_SLAVE = viper.GetString("POSTGRES_USERNAME_SLAVE")
	POSTGRES_PASSWORD_SLAVE = viper.GetString("POSTGRES_PASSWORD_SLAVE")
	POSTGRES_DATABASE_SLAVE = viper.GetString("POSTGRES_DATABASE_SLAVE")
	POSTGRES_SSL_MODE_SLAVE = viper.GetString("POSTGRES_SSL_MODE_SLAVE")

	ELASTICSEARCH_CLUSTER = viper.GetString("ELASTICSEARCH_CLUSTER")
	ELASTICSEARCH_USERNAME = viper.GetString("ELASTICSEARCH_USERNAME")
	ELASTICSEARCH_PASSWORD = viper.GetString("ELASTICSEARCH_PASSWORD")
	ELASTICSEARCH_PRODUCT_INDEX = viper.GetString("ELASTICSEARCH_PRODUCT_INDEX")
}
