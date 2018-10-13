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
		atlas, err := db.Open(indexDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		stime := time.Now()
		notes, err := atlas.QueryString(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found %v notes in %v\n", len(notes), time.Since(stime))
		eidx := len(notes)
		if eidx > 3  {
			eidx = 3
		}
		for _, note := range notes[:eidx] {
			if full {
				fmt.Printf("%v:\n\033[96m%v\033[39;49m\n", note.Title, note.Fragments)
			} else {
				fmt.Println(note.ID)
			}
		}
	},
}

var full = false
func init() {
	rootCmd.AddCommand(exqCmd)
	rootCmd.PersistentFlags().BoolVar(&full, "full", false, "include full results")

}
