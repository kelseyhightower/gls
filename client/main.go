package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/gls"
)

var (
	tlsCACert string
	tlsCert   string
	tlsKey    string
)

func homedir() string {
	return os.Getenv("HOME")
}

func init() {
	flag.StringVar(&tlsCACert, "tlscacert", filepath.Join(homedir(), ".gls/ca.pem"), "path to TLS CA cert")
	flag.StringVar(&tlsCert, "tlscert", filepath.Join(homedir(), ".gls/cert.pem"), "path to TLS cert")
	flag.StringVar(&tlsKey, "tlskey", filepath.Join(homedir(), ".gls/key.pem"), "path to TLS key")
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
