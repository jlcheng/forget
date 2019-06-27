package app

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/db/files"
	"github.com/jlcheng/forget/trace"
	"github.com/jlcheng/forget/watcher"
	"github.com/pkg/errors"
	"os"
	"time"
)

const RebuildBatchSize = 1000

// StartAtlasServer runs a server in a goroutine.
//
// Runs an Atlas server on the given port. The server will poll the dataDirs according to the pollInterval.
// Any new or modified files will be reindexed. The indexDir will be used to store or reuse the search index.
//
// If rebuild is true, the indexDir will be deleted and recreated.
func StartAtlasServer(port int, indexDir string, dataDirs []string, rebuild bool, pollInterval time.Duration) error {
	// if rebuild, delete the indexDir and recreate it. then reindex all in dataDirs
	if rebuild {
		err := rebuildIndexAndCloseAtlas(indexDir, dataDirs, RebuildBatchSize)
		if err != nil {
			return err
		}
	}

	// Start the Atlas server on port
	server := watcher.NewWatcherFacade()
	err := server.Listen(port, indexDir, dataDirs, pollInterval)
	if err != nil {
		return err
	}
	defer server.Close()

	return nil
}

func rebuildIndexAndCloseAtlas(indexDir string, dataDirs []string, batchSize int) error {
	stime := time.Now()
	err := os.RemoveAll(indexDir)
	if err != nil {
		return errors.Wrap(err, "Cannot delete IndexDir")
	}

	atlas, err := db.Open(indexDir, batchSize)
	if err != nil {
		return errors.Wrap(err, "Cannot open IndexDir")
	}
	defer atlas.CloseQuietly()
	err = files.RebuildIndex(atlas, dataDirs)
	if err != nil {
		return err
	}
	err = atlas.Flush()
	if err != nil {
		return err
	}
	trace.Info(fmt.Sprintf("Reindex complete (%v)", time.Since(stime)))
	return nil
}
