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
	ConfigParam   = "config"
	DataDirsParam = "dataDirs"
	HostParam     = "host"
	IndexDirParam = "indexDir"
	LogLevelParam = "logLevel"
	PortParam     = "port"
	PprofParam    = "pprof"
)

var ConfigFile = ""

func Config() string {
	return ConfigFile
}

func DataDirs() []string {
	return viper.GetStringSlice(DataDirsParam)
}

func Host() string {
	return viper.GetString(HostParam)
}

func IndexDir() string {
	return viper.GetString(IndexDirParam)
}

func Port() int {
	return viper.GetInt(PortParam)
}

func PprofEnabled() bool {
	return viper.GetBool(PprofParam)
}

func SetTraceLevel() {
	switch strings.ToUpper(viper.GetString(LogLevelParam)) {
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
	flags.StringVar(&ConfigFile, ConfigParam, "", "config file (default is $HOME/.forget.toml)")
	flags.StringSlice(DataDirsParam, make([]string, 0), "data directories")
	flags.String(HostParam, "localhost", "hostname of the 4gtsvr")
	flags.StringP(IndexDirParam, "i", "", "path to the index directory")
	flags.StringP(LogLevelParam, "L", "None", "log level: NONE, DEBUG, or WARN")
	flags.IntP(PortParam, "p", 8181, "rpc port")
	flags.Bool(PprofParam, false, "enable pprof server")
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

	SetTraceLevel()
	return nil
}
