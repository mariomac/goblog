// Package fs contains the File System tools
package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mariomac/goblog/src/logr"
	"github.com/sirupsen/logrus"
)

var log = logr.Get()

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
const gracePeriod = 5 * time.Second

func NotifyChanges(folder string, listener func() error) error {
	nlog := log.WithField("folder", folder)
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
		nlog.WithFields(logrus.Fields{
			"path":        path,
			"dirEntry":    d.Name(),
			"receivedErr": err,
		}).Debug("adding directory entry to notifier")
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

func processEvents(watcher *fsnotify.Watcher, nlog *logrus.Entry, listener func() error) {
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
				nlog.WithFields(logrus.Fields{
					"event":       event,
					"gracePeriod": gracePeriod,
				}).Info("reloading blog after grace period")
				time.AfterFunc(gracePeriod, func() {
					if err := listener(); err != nil {
						nlog.WithError(err).Error("couldn't reload blog")
					}
					// accept events again
					atomic.StoreInt64(&ignoreEvents, 0)
				})
			} else {
				nlog.WithFields(logrus.Fields{
					"event":       event,
					"gracePeriod": gracePeriod,
				}).Debug("still in grace period. Ignoring event")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				nlog.Warn("errors' channel closed. Exiting file watcher." +
					" Your blog won't be automatically updated after file changes")
				return
			}
			nlog.WithError(err).Error("error during file watching")
		}
	}
}
