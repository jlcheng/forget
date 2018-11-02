package cmd

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&CliCfg.CfgFile, "config", "", "config file (default is $HOME/.forget.yaml)")

	// Our own custom flags
	rootCmd.PersistentFlags().StringP("indexDir", "i", "", "path to the index directory")
	rootCmd.PersistentFlags().StringP("logLevel", "L", "None", "log level: NONE, DEBUG, or WARN")

	viper.BindPFlag(cli.INDEX_DIR, rootCmd.PersistentFlags().Lookup(cli.INDEX_DIR))
	viper.BindPFlag(cli.LOG_LEVEL, rootCmd.PersistentFlags().Lookup(cli.LOG_LEVEL))
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
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
