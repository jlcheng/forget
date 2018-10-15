package cmd

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var exqCmd = &cobra.Command{
	Use:   "exq",
	Short: "Query the index",
	Long: `Runs a query against the index`,
	Run: func(cmd *cobra.Command, args []string) {
		//os.Stdout.WriteString("\033")
		fmt.Println("exq called with:", strings.Join(args, " "))
		atlas, err := db.Open(indexDir, 1024)
		if err != nil {
			fmt.Println(err)
			return
		}
		docCount, err := atlas.GetDocCount()
		fmt.Println("atlas size:", docCount)
		stime := time.Now()
		notes, err := atlas.QueryString(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found %v notes in %v\n", len(notes), time.Since(stime))
		eidx := len(notes)
		if limit != 0 && eidx >= limit {
			eidx = limit
		}
		for _, note := range notes[:eidx] {
			if full {
				fmt.Printf("%v:\n\033[96m%v\033[39;49m\n", note.Title, note.Fragments["Body"])
			} else {
				fmt.Println(note.ID)
			}
		}
	},
}

var full = false
var limit = 0
func init() {
	rootCmd.AddCommand(exqCmd)
	rootCmd.PersistentFlags().BoolVar(&full, "full", false, "include full results")
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "limit results")
}
