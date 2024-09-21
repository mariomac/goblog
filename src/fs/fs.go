// Package fs contains the File System tools
package fs

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/mariomac/goblog/src/logr"
)

// Search returns the paths of all the files contained in the folder and subfolders of the path whose
// file name matches with the regular expression. If the regular expression is nil, it returns
// all the files
// The paths are returned in alphabetical order.
// It excludes the directories.
func Search(folder string, regexp *regexp.Regexp) ([]string, error) {
	var paths []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		// If there is any error, just ignores the file
		if err == nil && !info.IsDir() && (regexp == nil || regexp.MatchString(info.Name())) {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

// it will wait the grace period before notifying the listener. More events during that period
// will be ignored
// TODO: make configurable
const gracePeriod = 2 * time.Second

func NotifyChanges(folder string, listener func() error) error {
	nlog := logr.Get().With("folder", folder)
	nlog.Info("start file changes notifier")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// TODO: if folders are added or removed, we should add/remove them from watcher too
	if err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}
		nlog.Debug("adding directory entry to notifier",
			"path", path,
			"dirEntry", d.Name(),
			"receivedErr", err,
		)
		if err := watcher.Add(path); err != nil {
			return fmt.Errorf("adding folder %s: %w", path, err)
		}
		return nil
	}); err != nil {
		return err
	}
	go processEvents(watcher, nlog, listener)
	return nil
}

func processEvents(watcher *fsnotify.Watcher, nlog *slog.Logger, listener func() error) {
	defer watcher.Close()
	ignoreEvents := int64(0)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				nlog.Warn("events' channel closed. Exiting file watcher." +
					" Your blog won't be automatically updated after file changes")
				return
			}
			if atomic.CompareAndSwapInt64(&ignoreEvents, 0, 1) {
				nlog.Info("reloading blog after grace period",
					"event", event,
					"gracePeriod", gracePeriod,
				)
				time.AfterFunc(gracePeriod, func() {
					if err := listener(); err != nil {
						nlog.Error("couldn't reload blog", "error", err)
					}
					// accept events again
					atomic.StoreInt64(&ignoreEvents, 0)
				})
			} else {
				nlog.Debug("still in grace period. Ignoring event",
					"event", event,
					"gracePeriod", gracePeriod,
				)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				nlog.Warn("errors' channel closed. Exiting file watcher." +
					" Your blog won't be automatically updated after file changes")
				return
			}
			nlog.Error("error during file watching", "error", err)
		}
	}
}
