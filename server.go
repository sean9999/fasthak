package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// set up TLS connection
func ListenAndServeTLSKeyPair(addr string, cert tls.Certificate, handler http.Handler) error {
	if addr == "" {
		return errors.New("invalid address string")
	}
	server := &http.Server{Addr: addr, Handler: handler}
	config := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	config.NextProtos = []string{"http/1.1"}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = cert
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	return server.Serve(tlsListener)
}

func injectHeadersForStaticFiles(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, max-age=1")
		fs.ServeHTTP(w, r)
	}
}

func serve(staticDir string, port int, privKeyPath string, pubKeyPath string, broker *Broker) error {
	//	get TLS keys. panic if they don't exist
	pubKeyMaterial, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		panic("Could not find " + pubKeyPath)
	}
	privKeyMaterial, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		panic("Could not find " + privKeyPath)
	}
	//	start web server
	cert, err := tls.X509KeyPair(pubKeyMaterial, privKeyMaterial)
	if err != nil {
		log.Fatalln(err)
	}
	mux := http.NewServeMux()

	//	static files
	staticFileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/", injectHeadersForStaticFiles(staticFileServer))

	//	.hak/js/*
	mux.Handle(hakPrefix+"/js/", hakHandler())

	//	push sseBroker events to /./hak/fs/sse
	mux.Handle(path.Join(hakPrefix, ssePath), broker)

	portString := fmt.Sprintf("%s%d", ":", port)
	fmt.Printf("listening on port %d\n", port)
	err = ListenAndServeTLSKeyPair(portString, cert, mux)

	return err

}

// embedded client-side /.hak/js/*
func hakHandler() http.Handler {
	fsys := fs.FS(frontend)
	hakFiles, _ := fs.Sub(fsys, "frontend")
	return http.StripPrefix(hakPrefix+"/js/", http.FileServer(http.FS(hakFiles)))
}
