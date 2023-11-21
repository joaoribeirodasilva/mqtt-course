package main

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

type TlsConfig struct {
	caRootFile string
	certFile   string
	keyFile    string
	insecure   bool
	Config     *tls.Config
}

func NewTlsConfig(certFile string, keyFile string, caRootFile string, insecure bool) *TlsConfig {

	t := &TlsConfig{}

	t.caRootFile = caRootFile
	t.certFile = certFile
	t.keyFile = keyFile
	t.insecure = insecure

	return t
}

func (t *TlsConfig) Create() error {

	// create a new certificate pool (root(if set) + public certificate + private certificate)
	caCertPool := x509.NewCertPool()

	// if a CA root certificate file is set
	if t.caRootFile != "" {

		// read CA root file
		caRoot, err := os.ReadFile(t.caRootFile)
		if err != nil {
			return err
		}

		// append the CA root certificate to the certificate pool
		// and check if it is valid
		if ok := caCertPool.AppendCertsFromPEM(caRoot); !ok {
			return err
		}
	}

	// load the public and private certificate files
	certs, err := tls.LoadX509KeyPair(t.certFile, t.keyFile)
	if err != nil {
		return err
	}

	// instantiate the TLS configuration
	t.Config = &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: t.insecure,
		Certificates:       []tls.Certificate{certs},
	}

	return nil
}
