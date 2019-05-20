package main

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/rpc"
	"github.com/jlcheng/forget/txtio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

// Keys used to store and lookup values through the Viper API
const (
	nameHost = "host"
	namePort = "port"
)

// This must be global so it can be shared between the Cobra and Viper APIs
var config string = ""

// This must be global so we can bind values parsed by Cobra to the Viper API
var rootCmd = &cobra.Command{
	Use:   "4gt",
	Short: "Query the 4gt server",
	Long:  "Query the 4gt server",
	Run:   cobraMain,
}

func cobraMain(cmd *cobra.Command, args []string) {
	cli.SetTraceLevel()
	host := viper.GetString(nameHost)
	port := viper.GetInt(namePort)

	qterms := make([]string, len(args))
	for idx := range args {
		qterms[idx] = "+Body:" + args[idx]
	}
	SendQuery(host, port, strings.Join(qterms, " "))
}

// Loads configuration files
func cobraInit() {
	if config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config)
	} else {
		// Search config in home directory with name ".forget" (without extension).
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".forget")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error loading config file.", err)

	}
	if viper.ConfigFileUsed() != "" {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Bind flags from the command line to the viper framework
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatal(err)
	}
}

// TODO: Move this into the main directory
func SendQuery(host string, port int, query string) {
	stime := time.Now()
	atlasResponse, err := rpc.Request(host, port, query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %v lines in %v\n", len(atlasResponse.ResultEntries), time.Since(stime))
	for _, entry := range atlasResponse.ResultEntries {
		fmt.Println(txtio.AnsiFmt(entry))
	}
}

func main() {
	rootCmd.Flags().StringVar(&config, "config", "", "config file (default: $HOME/.forget.toml)")
	rootCmd.Flags().IntP(namePort, "p", 8181, "rpc port")
	rootCmd.Flags().StringP(nameHost, "H", "localhost", "rpc host")

	cobra.OnInitialize(cobraInit)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
