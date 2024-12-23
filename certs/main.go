package certs

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
)

/**
 *	We compress these certs as a courtesy to backloop.dev and Git Guardian.
 *	Automated scanners are likely to flag unobfuscated private key material.
 **/

//go:embed *.gz
var folder embed.FS

func decompress(gz fs.File, err error) ([]byte, error) {
	if err != nil {
		return nil, fmt.Errorf("could not decompress. %w", err)
	}
	zr, _ := gzip.NewReader(gz)
	buff := new(bytes.Buffer)
	buff.ReadFrom(zr)
	return buff.Bytes(), nil
}

func getFile(filesystem fs.FS, name string) (fs.File, error) {
	f, err := filesystem.Open(name)
	if err != nil {
		return nil, fmt.Errorf("could not get file. %w", err)
	}
	return f, nil
}

func KeyPair() (tls.Certificate, error) {
	var noCert tls.Certificate
	cert, err := decompress(getFile(folder, "backloop.dev-ca.crt"))
	if err != nil {
		return noCert, fmt.Errorf("could not create key-pair. %w", err)
	}
	key, err := decompress(getFile(folder, "backloop.dev-key.pem"))
	if err != nil {
		return noCert, fmt.Errorf("could not create key-pair. %w", err)
	}
	return tls.X509KeyPair(cert, key)
}
