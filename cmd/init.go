package cmd

import (
	"context"
	stdlog "log"
	"reflect"

	"github.com/go-playground/validator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/proemergotech/log-go"
)

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

		err := viper.BindEnv(name)
		if err != nil {
			stdlog.Printf("config error: " + err.Error())
			hasErrors = true
			continue
		}

		def := fieldType.Tag.Get("default")
		if def != "" {
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
	err := validate.Struct(cfg)
	if err != nil {
		log.Panic(context.Background(), "invalid configuration", "error", err)
	}
}
