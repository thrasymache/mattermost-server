package config

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/mlog"
)

// watcher monitors a file for changes
type watcher struct {
	emitter

	fsWatcher *fsnotify.Watcher
	close     chan struct{}
	closed    chan struct{}
}

// newWatcher creates a new instance of watcher to monitor for file changes.
func newWatcher(path string, callback func()) (*watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create watcher for %s", path)
	}

	path = filepath.Clean(path)

	// Watch the entire containing directory.
	configDir, _ := filepath.Split(path)
	fsWatcher.Add(configDir)

	ret := &watcher{
		fsWatcher: fsWatcher,
		close:     make(chan struct{}),
		closed:    make(chan struct{}),
	}

	go func() {
		defer close(ret.closed)
		defer fsWatcher.Close()

		for {
			select {
			case event := <-fsWatcher.Events:
				// We only care about the given file.
				if filepath.Clean(event.Name) == path {
					if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
						mlog.Info("Config file watcher detected a change", mlog.String("path", path))
						callback()
					}
				}
			case err := <-fsWatcher.Errors:
				mlog.Error("Failed while watching config file", mlog.String("path", path), mlog.Err(err))
			case <-ret.close:
				return
			}
		}
	}()

	return ret, nil
}

func (w *watcher) Close() error {
	close(w.close)
	<-w.closed

	return nil
}
