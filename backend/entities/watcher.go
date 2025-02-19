package entities

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

type FileWatcher interface {
	Start() error
	Close() error
}

type fsNotifyWatcher struct {
	watcher      *fsnotify.Watcher
	path         string
	onFileChange func()
}

func NewFsNotifyWatcher(path string, onFileChange func()) (FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &fsNotifyWatcher{
		watcher:      watcher,
		path:         path,
		onFileChange: onFileChange,
	}, nil
}

func (w *fsNotifyWatcher) Start() error {
	if err := w.watcher.Add(w.path); err != nil {
		log.Error().Err(err).Msgf("Error watching file: %s", w.path)
		return w.Close()
	}

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					w.onFileChange()
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Error().Err(err).Msg("Watcher error")
			}
		}
	}()

	return nil
}

func (w *fsNotifyWatcher) Close() error {
	return w.watcher.Close()
}
