//handles the configuration of the applications. Yaml files are mapped with the struct

package config

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	"github.com/zalando/gin-contrib/ginoauth2"
	"golang.org/x/oauth2"
)

// User contain fields for
// OAUTH2 implementation
type User struct {
	Username string
	Fullname string
	Role     string
	Group    string
}

// Config contain fields for config yaml file
type Config struct {
	Endpoints      map[string]string
	AllowedUsers   []User
	Security       map[string]string
	Documentation  map[string]string
	Ltmdevicenames map[string]string
}

// ConfigError contain fields of config error handling
//created a struct just for future usage
type ConfigError struct {
	Message string
}

// Error return config error
// viper can not parse yaml file
func (e *ConfigError) Error() string {
	return fmt.Sprintf(e.Message)
}

// ConfigInit initiliaze configuration file
func ConfigInit(filename string) (*Config, *ConfigError) {
	viper.SetConfigType("YAML")
	f, err := os.Open(filename)
	if err != nil {
		return nil, &ConfigError{"could not read configuration files."}
	}
	err = viper.ReadConfig(f)
	if err != nil {
		return nil, &ConfigError{"configuration format is not correct."}
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		glog.Errorf("Cannot read configuration. Reason: %s", err)
		return nil, &ConfigError{"cannot read configuration, something must be wrong."}
	}

	return &config, nil
}

// LoadAuthConf extract necessary oAuth2 information
func LoadAuthConf(config *Config) ([]ginoauth2.AccessTuple, []ginoauth2.AccessTuple, *oauth2.Endpoint, *ConfigError) {
	var rootUsers = []ginoauth2.AccessTuple{}
	var emergencyUsers = []ginoauth2.AccessTuple{}
	var endpoint = oauth2.Endpoint{}
	for _, user := range config.AllowedUsers {
		username := user.Username
		fullname := user.Fullname
		role := user.Role
		if username == "" || fullname == "" || role == "" {
			return nil, nil, nil, &ConfigError{"configuration is invalid. TokenRUL or AuthURL are missing"}
		}
		u := ginoauth2.AccessTuple{role, username, fullname}
		if user.Group == "root" {
			rootUsers = append(rootUsers, u)
		}
		if user.Group == "emergency" || user.Group == "root" {
			emergencyUsers = append(emergencyUsers, u)
		}
	}
	authURL := config.Endpoints["AuthURL"]
	tokenURL := config.Endpoints["TokenURL"]
	if authURL == "" || tokenURL == "" {
		return nil, nil, nil, &ConfigError{"configuration is invalid. TokenURL or AuthURL are missing"}
	}
	endpoint = oauth2.Endpoint{authURL, tokenURL}
	return rootUsers, emergencyUsers, &endpoint, nil
}

// LoadConfig load config file
func LoadConfig() *Config {
	var err *ConfigError
	conf, err := ConfigInit("config.yaml")
	if err != nil {
		glog.Errorf("Cannot load configuration. Reason: %s", err.Message)
		panic("Cannot load configuration for Baboon. Exiting.")
	}
	return conf
}
