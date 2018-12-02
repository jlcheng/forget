package cli

import (
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/viper"
	"strings"
)


const (
	INDEX_DIR = "indexDir"
	LOG_LEVEL = "logLevel"
)


type CLIConfig struct {
	CfgFile string
	IndexDir string
	DataDirs []string
}

func (c *CLIConfig) GetIndexDir() string {
	return viper.GetString(INDEX_DIR)
}

func (c *CLIConfig) SetTraceLevel() {
	switch strings.ToUpper(viper.GetString(LOG_LEVEL)) {
	case "DEBUG":
		trace.Level = trace.LOG_DEBUG
	case "WARN":
		trace.Level = trace.LOG_WARN
	default:
		trace.Level = trace.LOG_NONE
	}
}