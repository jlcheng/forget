package cli

import (
	"fmt"
	"github.com/jlcheng/forget/trace"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	CONFIG        = "config"
	DATA_DIRS     = "dataDirs"
	HOST          = "host"
	INDEX_DIR     = "indexDir"
	LOG_LEVEL     = "logLevel"
	PORT          = "port"
	PPROF_ENABLED = "pprof"
)

var ConfigFile = ""

func Config() string {
	return ConfigFile
}

func DataDirs() []string {
	return viper.GetStringSlice(DATA_DIRS)
}

func Host() string {
	return viper.GetString(HOST)
}

func IndexDir() string {
	return viper.GetString(INDEX_DIR)
}

func Port() int {
	return viper.GetInt(PORT)
}

func PprofEnabled() bool {
	return viper.GetBool(PPROF_ENABLED)
}

func SetTraceLevel() {
	switch strings.ToUpper(viper.GetString(LOG_LEVEL)) {
	case "DEBUG":
		trace.Level = trace.LOG_DEBUG
	case "WARN":
		trace.Level = trace.LOG_WARN
	default:
		trace.Level = trace.LOG_NONE
	}
}

// ConfigureFlagSet configures the given *pflag.FlagSet for parsing
func ConfigureFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&ConfigFile, CONFIG, "", "config file (default is $HOME/.forget.toml)")
	flags.StringSlice(DATA_DIRS, make([]string, 0), "data directories")
	flags.String(HOST, "localhost", "hostname of the 4gtsvr")
	flags.StringP(INDEX_DIR, "i", "", "path to the index directory")
	flags.StringP(LOG_LEVEL, "L", "None", "log level: NONE, DEBUG, or WARN")
	flags.IntP(PORT, "p", 8181, "rpc port")
	flags.Bool(PPROF_ENABLED, false, "enable pprof server")
}

// ProcessParsedFlagSet uses a parsed *pflag.FlagSet to kick-off processing of the config file
func ProcessParsedFlagSet(flags *pflag.FlagSet) error {
	if Config() != "" {
		// Use config file from the flag.
		viper.SetConfigFile(Config())
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return errors.WithStack(err)
		}

		// Search config in home directory with name ".forget" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".forget")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return errors.WithStack(err)

	}
	if viper.ConfigFileUsed() != "" {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Combine parsed flags with contents of the configuration file
	if err := viper.BindPFlags(flags); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
