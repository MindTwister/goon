package goon

import (
	"io/ioutil"
	"os"
	"path"
	"time"
)

func getAllFiles(dirname string) (map[string]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	out := make(map[string]os.FileInfo)
	for _, file := range files {
		if file.IsDir() {
			tmp, _ := getAllFiles(path.Join(dirname, file.Name()))
			for path, f := range tmp {
				out[path] = f
			}
		} else {
			out[path.Join(dirname, file.Name())] = file
		}
	}
	return out, nil
}

func Watch(dirname string, interval int) (changed chan struct{}) {
	changed = make(chan struct{})
	go (func() {
		files := make(map[string]os.FileInfo)
		for {
			isModified := false
			tmp, err := getAllFiles(dirname)
			if err != nil {
				close(changed)
			}
			for name, stat := range tmp {
				file, ok := files[name]
				if !ok {
					isModified = true
					break
				}
				if file.ModTime().Unix() < stat.ModTime().Unix() {
					isModified = true
					break
				}
			}
			if len(tmp) != len(files) {
				isModified = true
			}
			if isModified {
				changed <- struct{}{}
			}
			files = tmp
			time.Sleep(time.Millisecond * time.Duration(interval))
		}
	})()
	return changed
}
