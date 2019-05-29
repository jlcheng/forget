package subcmd

import (
	"github.com/jlcheng/forget/app"
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/cobra"
	"os"
)

var qcCmd = &cobra.Command{
	Use:   "qc",
	Short: "Queries the 4gt server",
	Long:  "Queries a 4gt server and display the results using grep-like output",
	Run: func(_*cobra.Command, args []string) {
		if err := app.GrepClient(args); err != nil {
			trace.PrintStackTrace(err)
			os.Exit(1)
		}
	},
}

func InitExqc() {
	rootCmd.AddCommand(qcCmd)
}
