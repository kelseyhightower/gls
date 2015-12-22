package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/kelseyhightower/gls"
)

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	files := make(gls.Files, 0)
	err = client.Call("Ls.Ls", os.Args[1], &files)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Printf("%s %10d %s %s\n", f.Mode, f.Size, f.ModTime, f.Name)
	}
}
