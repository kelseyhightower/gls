package gls

import (
	"os"
	"path/filepath"
)

type Files []File

type File struct {
	Name    string
	Size    int64
	Mode    string
	ModTime string
}

type Ls struct{}

func (ls *Ls) Ls(path *string, files *Files) error {
	root := *path
	err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file := File{
			info.Name(),
			info.Size(),
			info.Mode().String(),
			info.ModTime().Format("Jan _2 15:04"),
		}
		*files = append(*files, file)
		if info.IsDir() && path != root {
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
