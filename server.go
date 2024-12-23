package main

import (
	"crypto/tls"
	"errors"
	"io/fs"
	"net"
	"net/http"
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
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{cert},
		ServerName:   domain,
	}
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		TLSConfig:    &tls.Config{},
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	//return srv.ListenAndServeTLS(certFile string, keyFile string)
	tlsListener := tls.NewListener(ln, cfg)
	return srv.Serve(tlsListener)
}

func injectHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, max-age=1")
		fs.ServeHTTP(w, r)
	}
}

// embedded client-side /.hak/js/*
func hakHandler() http.Handler {
	fsys := fs.FS(frontend)
	hakFiles, _ := fs.Sub(fsys, "frontend")
	return http.StripPrefix(hakPrefix+"/js/", http.FileServer(http.FS(hakFiles)))
}
