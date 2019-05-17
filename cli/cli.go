package cli

import (
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/viper"
	"strings"
)

const (
	INDEX_DIR     = "indexDir"
	DATA_DIRS     = "dataDirs"
	LOG_LEVEL     = "logLevel"
	PPROF_ENABLED = "pprof"
)

type CLIConfig struct {
	CfgFile string
}

func (c *CLIConfig) GetIndexDir() string {
	return viper.GetString(INDEX_DIR)
}

func (c *CLIConfig) GetDataDirs() []string {
	return viper.GetStringSlice(DATA_DIRS)
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
