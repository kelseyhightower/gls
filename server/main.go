package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/kelseyhightower/gls"
)

func main() {
	log.Println("Starting glsd..")
	ls := new(gls.Ls)
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
