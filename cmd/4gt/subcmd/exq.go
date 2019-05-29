package subcmd

import (
	"github.com/jlcheng/forget/app"
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/cobra"
	"os"
)

var exqCmd = &cobra.Command{
	Use:   "exq",
	Short:  "Experimental query client",
	Long:  "Experimental query client",
	Run: func(_ *cobra.Command, args []string) {
		if err := app.SnippetClient(args); err != nil {
			trace.PrintStackTrace(err)
			os.Exit(1)
		}
	},
}

var exqLimit = 0

func InitExq() {
	rootCmd.AddCommand(exqCmd)
	rootCmd.Flags().IntVarP(&exqLimit, "limit", "l", 0, "limit results")
}
