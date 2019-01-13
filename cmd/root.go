package cmd

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var CliCfg = cli.CLIConfig{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "4gt",
	Short: "A personal information management system",
	Long: `Forget is a CLI program to index and find information for the absent minded.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&CliCfg.CfgFile, "config", "", "config file (default is $HOME/.forget.toml)")

	pflags.StringP(cli.INDEX_DIR, "i", "", "path to the index directory")
	if err := viper.BindPFlag(cli.INDEX_DIR, rootCmd.PersistentFlags().Lookup(cli.INDEX_DIR)); err != nil {
		log.Fatal(err)
	}

	pflags.StringP(cli.LOG_LEVEL, "L", "None", "log level: NONE, DEBUG, or WARN")
	if err := viper.BindPFlag(cli.LOG_LEVEL, rootCmd.PersistentFlags().Lookup(cli.LOG_LEVEL)); err != nil {
		log.Fatal(err)
	}

	pflags.StringSlice(cli.DATA_DIRS, make([]string,0,0), "data directories")
	if err := viper.BindPFlag(cli.DATA_DIRS, rootCmd.PersistentFlags().Lookup(cli.DATA_DIRS)); err != nil {
		log.Fatal(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
