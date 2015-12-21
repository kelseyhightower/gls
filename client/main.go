package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"

	"github.com/kelseyhightower/gls"
)

var (
	tlsCACert string
	tlsCert   string
	tlsKey    string
)

func init() {
	flag.StringVar(&tlsCACert, "tlscacert", "", "path to TLS CA cert")
	flag.StringVar(&tlsCert, "tlscert", "", "path to TLS cert")
	flag.StringVar(&tlsKey, "tlskey", "", "path to TLS key")
}

func main() {
	flag.Parse()

	caCert, err := ioutil.ReadFile(tlsCACert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:8080", &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	})
	if err != nil {
		log.Fatal(err)
	}

	client := rpc.NewClient(conn)

	files := make(gls.Files, 0)
	err = client.Call("Ls.Ls", flag.Args()[0], &files)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Printf("%s %10d %s %s\n", f.Mode, f.Size, f.ModTime, f.Name)
	}
}
