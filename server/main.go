package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/gls"
)

type Ls struct{}

func (ls *Ls) Ls(path *string, files *gls.Files) error {
	root := *path
	err := filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file := gls.File{
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

func main() {
	log.Println("Starting glsd..")
	ls := new(Ls)
	rpc.Register(ls)
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		rpc.ServeConn(conn)
		conn.Close()
	}
}
