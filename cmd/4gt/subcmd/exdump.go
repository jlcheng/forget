package subcmd

import (
	"fmt"
	"github.com/jlcheng/forget/cli"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/trace"
	"github.com/spf13/cobra"
	"log"
)

// exdumpCmd represents the exdump command
var exdumpCmd = &cobra.Command{
	Use:   "exdump",
	Short: "Dumps the index into the specified directory.",
	Long:  "Dumps the index into the specified directory.",
	Run: func(cmd *cobra.Command, args []string) {
		cli.SetTraceLevel()
		fmt.Println("exdump called")
		atlas, err := db.Open(cli.IndexDir(), 16)
		if err != nil {
			log.Fatal(err)
		}
		if err = IterateDocuments(atlas); err != nil {
			log.Fatal(err)
		}
		atlas.Close()
	},
}

func InitExdump() {
	rootCmd.AddCommand(exdumpCmd)
}

// IterateDocuments
func IterateDocuments(atlas *db.Atlas) error {
	trace.Debug("IterateDocuments")
	docs, err := atlas.DumpAll()
	if err != nil {
		return err
	}
	for _, doc := range docs {
		trace.Debug(doc.PrettyString())
	}
	return nil
}
