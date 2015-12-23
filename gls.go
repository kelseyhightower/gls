package gls

type Files []File

type File struct {
	Name    string
	Size    int64
	Mode    string
	ModTime string
}
