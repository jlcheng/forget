package subcmd

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "4gt",
	Short: "A personal information management system",
	Long:  `Forget is a CLI program to index and find information for the absent-minded.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func InitCobra() {
	cobra.OnInitialize(cobraInit)

	flags := rootCmd.PersistentFlags()
	flags.StringVar(&cli.ConfigFile, "config", "", "config file (default is $HOME/.forget.toml)")

	flags.StringP(cli.INDEX_DIR, "i", "", "path to the index directory")
	flags.StringP(cli.LOG_LEVEL, "L", "None", "log level: NONE, DEBUG, or WARN")
	flags.StringSlice(cli.DATA_DIRS, make([]string, 0), "data directories")
	flags.IntP(cli.PORT, "p", 8181, "rpc port")

	InitExdump()
	InitExload()
	InitExqc()
	InitExq()
	InitExsvr()
}

// cobraInit reads in config file and ENV variables if set.
func cobraInit() {
	if err := cli.ProcessParsedFlagSet(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}

}
