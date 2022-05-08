package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/rjeczalik/notify"
)

// From https://golang.org/src/net/http/server.go
// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
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

// ListenAndServeTLSKeyPair start a server using in-memory TLS KeyPair
func ListenAndServeTLSKeyPair(addr string, cert tls.Certificate,
	handler http.Handler) error {

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

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)},
		config)

	return server.Serve(tlsListener)
}

func injectHeadersForStaticFiles(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, max-age=1")
		fs.ServeHTTP(w, r)
	}
}

func handler(f http.HandlerFunc) http.Handler {
	return f
}

func pushEvent(ei notify.EventInfo, w http.ResponseWriter) {
	/**
	 *	push a fileSystem Event through SSE
	 */

	//buf, err := notifyBuff(ei)

	ne := toNiceEvent(ei)
	log.Printf("%s - %s", ne.Event, ne.File)

	buf, err := niceEventToBuffer(ne)
	if err != nil {
		log.Panic(err)
		return
	}
	_, err2 := fmt.Fprintf(w, "event: %s\ndata: %s\nretry: 3000\n\n", "fs", buf.String())
	if err2 != nil {
		log.Printf("client disconnect - %s", err2)
	}
	w.(http.Flusher).Flush()
}
