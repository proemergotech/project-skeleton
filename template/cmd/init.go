package cmd

import (
	"context"
	stdlog "log"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gitlab.com/proemergotech/log-go/v3"
)

// initConfig reads in config file and ENV variables if set.
func initConfig(cfg interface{}) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			log.Panic(context.Background(), "specified config file does not exists", "config_file", cfgFile)
		}
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Panic(context.Background(), "unable to read config", "err", err)
		}
	}

	hasErrors := false
	val := reflect.ValueOf(cfg).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		name := fieldType.Tag.Get("mapstructure")
		if name == "" {
			stdlog.Printf("Config error: settings struct field " + fieldType.Name + " has no mapstructure tag")
			hasErrors = true
			continue
		}

		if err := viper.BindEnv(name); err != nil {
			stdlog.Printf("config error: " + err.Error())
			hasErrors = true
			continue
		}

		if def := fieldType.Tag.Get("default"); def != "" {
			viper.SetDefault(name, def)
		}
	}

	if hasErrors {
		log.Panic(context.Background(), "config error happened, check the log for details")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panic(context.Background(), "Unable to marshal config", "error", err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.Panic(context.Background(), "invalid configuration", "error", err)
	}
}
