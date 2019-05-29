package subcmd

import (
	"github.com/jlcheng/forget/cli"
	"github.com/spf13/cobra"
	"log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "4gt",
	Short: "A personal information management system",
	Long:  `Forget is a CLI program to index and find information for the absent-minded.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func InitRoot() {
	cobra.OnInitialize(func(){
		if err := cli.ProcessParsedFlagSet(rootCmd.PersistentFlags()); err != nil {
			log.Fatal(err)
		}

	})

	cli.ConfigureFlagSet(rootCmd.PersistentFlags())

	InitExdump()
	InitExqc()
	InitExq()
	InitSvr()
}
