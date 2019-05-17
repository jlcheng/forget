package main


import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/watcher"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"time"	
)

import _ "net/http/pprof"

var CliCfg = cli.CLIConfig{}

func main() {
	cobra.OnInitialize(cobraInit)

	var port int
	var duration int
	
	var rootCmd = &cobra.Command{
		Use:   "4gt",
		Short: "Starts the 4gt server",
		Long:  `Starts the 4gt server`,
		Run: func(cmd *cobra.Command, args []string) {
			CliCfg.SetTraceLevel()

			// turn on pprof if specified
			if viper.GetBool(cli.PPROF_ENABLED) {
				go func() {
					fmt.Println("pprof running on port 6060")
					log.Println(http.ListenAndServe("localhost:6060", nil))
				}()
			}
			

			runexsvr := watcher.NewWatcherFacade()
			defer runexsvr.Close()
			err := runexsvr.Listen(port, CliCfg.GetIndexDir(), CliCfg.GetDataDirs(), time.Second*time.Duration(duration))
			if err != nil {
				trace.PrintStackTrace(err)
				os.Exit(1)
			}
		},
	}
	pflags := rootCmd.PersistentFlags()
	
	pflags.IntVarP(&port, "port", "p", 8181, "rpc port")
	pflags.IntVarP(&duration, "duration", "t", 10, "seconds between polling the file system for changes")
	pflags.StringP(cli.LOG_LEVEL, "L", "None", "log level: NONE, DEBUG, or WARN")
	if err := viper.BindPFlag(cli.LOG_LEVEL, pflags.Lookup(cli.LOG_LEVEL)); err != nil {
		log.Fatal(err)
	}
	
	pflags.Bool(cli.PPROF_ENABLED, false, "if specified, turn on pprof on port 6060")
	if err := viper.BindPFlag(cli.PPROF_ENABLED, pflags.Lookup(cli.PPROF_ENABLED)); err != nil {
		log.Fatal(err)
	}
	
	rootCmd.Execute()
}

// Loads configuration files
func cobraInit() {
	if CliCfg.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CliCfg.CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".forget" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".forget")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error loading config file.", err)

	}
	if viper.ConfigFileUsed() != "" {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
