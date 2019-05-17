package main


import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/rpc"
	"github.com/jlcheng/forget/txtio"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strings"
	"time"	
)

var CliCfg = cli.CLIConfig{}

func main() {
	cobra.OnInitialize(cobraInit)

	var host = "localhost"
	var port = 8181
	
	var rootCmd = &cobra.Command{
		Use:   "4gt",
		Short: "Query the 4gt server",
		Long:  `Query the 4gt server`,
		Run: func(cmd *cobra.Command, args []string) {
			CliCfg.SetTraceLevel()

			qterms := make([]string, len(args))
			for idx := range args {
				qterms[idx] = "+" + args[idx]
			}
			stime := time.Now()
			atlasResponse, err := rpc.Request(host, port, strings.Join(args, " "))
			if err != nil {
				log.Fatal(err)
			}
			
			fmt.Printf("Found %v notes in %v\n", len(atlasResponse.ResultEntries), time.Since(stime))
			for _, entry := range atlasResponse.ResultEntries {
				fmt.Println(txtio.AnsiFmt(entry))
			}
			return			
		},
	}
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8181, "rpc port")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "localhost", "rpc host")
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

	// turn on pprof if specified
	if viper.GetBool(cli.PPROF_ENABLED) {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}
