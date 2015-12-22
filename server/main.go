package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
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

func tlsPath(dir, file string) string {
	return filepath.Join(os.Getenv("HOME"), dir, file)
}

func init() {
	flag.StringVar(&tlsCACert, "tlscacert", tlsPath(".glsd", "ca.pem"), "path to TLS CA cert")
	flag.StringVar(&tlsCert, "tlscert", tlsPath(".glsd", "cert.pem"), "path to TLS cert")
	flag.StringVar(&tlsKey, "tlskey", tlsPath(".glsd", "key.pem"), "path to TLS key")
}

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
	flag.Parse()
	log.Println("Starting glsd..")

	ls := new(Ls)
	rpc.Register(ls)

	caCert, err := ioutil.ReadFile(tlsCACert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	serverCert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}
	tlsConfig.BuildNameToCertificate()

	l, err := tls.Listen("tcp", "0.0.0.0:8080", tlsConfig)
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
