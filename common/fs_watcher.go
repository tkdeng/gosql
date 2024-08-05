package common

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// A watcher instance for the `FS.FSWatcher` method
type FSWatcher struct {
	watcherList *map[string]*watcherObj
	mu          sync.Mutex
	size        *uint

	// when a file changes
	//
	// @path: the file path the change happened to
	//
	// @op: the change operation
	OnFileChange func(path string, op string)

	// when a directory is added
	//
	// @path: the file path the change happened to
	//
	// @op: the change operation
	//
	// return false to prevent that directory from being watched
	OnDirAdd func(path string, op string) (addWatcher bool)

	// when a file or directory is removed
	//
	// @path: the file path the change happened to
	//
	// @op: the change operation
	//
	// return false to prevent that directory from no longer being watched
	OnRemove func(path string, op string) (removeWatcher bool)

	// every time something happens
	//
	// @path: the file path the change happened to
	//
	// @op: the change operation
	OnAny func(path string, op string)
}

type watcherObj struct {
	watcher *fsnotify.Watcher
	close   *bool
}

// FileWatcher creates a new file watcher
func FileWatcher() *FSWatcher {
	size := uint(0)
	return &FSWatcher{watcherList: &map[string]*watcherObj{}, size: &size}
}

// WatchDir watches the files in a directory and its subdirectories for changes
//
// @nosub: do not watch sub directories
func (fw *FSWatcher) WatchDir(root string, nosub ...bool) error {
	var err error
	if root, err = filepath.Abs(root); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	runClose := false

	fw.mu.Lock()
	(*fw.watcherList)[root] = &watcherObj{watcher: watcher, close: &runClose}
	*fw.size++
	fw.mu.Unlock()

	go func() {
		defer watcher.Close()
		for {
			if runClose {
				break
			}

			if event, ok := <-watcher.Events; ok {
				filePath := event.Name

				stat, err := os.Stat(filePath)
				if err != nil {
					if fw.OnRemove == nil || fw.OnRemove(filePath, event.Op.String()) {
						watcher.Remove(filePath)
					}
				} else if stat.IsDir() {
					if fw.OnDirAdd == nil || fw.OnDirAdd(filePath, event.Op.String()) {
						watcher.Add(filePath)
					}
				} else {
					if fw.OnFileChange != nil {
						fw.OnFileChange(filePath, event.Op.String())
					}
				}

				if fw.OnAny != nil {
					fw.OnAny(filePath, event.Op.String())
				}
			}
		}
	}()

	err = watcher.Add(root)
	if err != nil {
		return err
	}

	if len(nosub) == 0 || nosub[0] == false {
		fw.watchDirSub(watcher, root)
	}

	return nil
}

func (fw *FSWatcher) watchDirSub(watcher *fsnotify.Watcher, dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			if path, err := JoinPath(dir, file.Name()); err == nil {
				watcher.Add(path)
				fw.watchDirSub(watcher, path)
			}
		}
	}
}

// CloseWatcher will close the watcher by the root name you used
//
// @root pass a file path for a specific watcher or "*" for all watchers that exist
func (fw *FSWatcher) CloseWatcher(root string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if root == "" || root == "*" {
		for r, w := range *fw.watcherList {
			*w.close = true
			delete(*fw.watcherList, r)
			*fw.size--
		}
	} else {
		var err error
		if root, err = filepath.Abs(root); err != nil {
			return err
		}

		if w, ok := (*fw.watcherList)[root]; ok {
			*w.close = true
			delete(*fw.watcherList, root)
			*fw.size--
		}
	}

	return nil
}

// Wait for all Watchers to close
func (fw *FSWatcher) Wait() {
	for *fw.size != 0 {
		time.Sleep(100 * time.Millisecond)
	}
}
