package subcmd

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/txtio"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	BATCH_SIZE = 1024
)

var exqCmd = &cobra.Command{
	Use:   "exq",
	Short: "Query the index",
	Long:  `Runs a query against the index`,
	Run: func(cmd *cobra.Command, args []string) {
		CliCfg.SetTraceLevel()

		atlas, err := db.Open(CliCfg.GetIndexDir(), BATCH_SIZE)
		if err != nil {
			fmt.Println(err)
			return
		}
		docCount, err := atlas.GetDocCount()
		if err != nil {
			fmt.Println(err)
			return
		}
		trace.Debug("atlas size:", docCount)
		stime := time.Now()
		atlasResponse := atlas.QueryForResponse(strings.Join(args, " "))
		fmt.Printf("found %v notes in %v\n", len(atlasResponse.ResultEntries), time.Since(stime))
		eidx := len(atlasResponse.ResultEntries)
		if limit != 0 && eidx >= limit {
			eidx = limit
		}
		for _, entry := range atlasResponse.ResultEntries[:eidx] {
			fmt.Println(txtio.AnsiFmt(entry))
		}
	},
}

var full = false
var limit = 0

func InitExq() {
	rootCmd.AddCommand(exqCmd)
	rootCmd.PersistentFlags().BoolVar(&full, "full", false, "include full results")
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "limit results")
}