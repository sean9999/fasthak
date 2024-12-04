package main

import (
	"crypto/tls"
	"encoding/base64"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func AnchorCert() *tls.Config {
	acmeKey, _ := base64.RawURLEncoding.DecodeString(os.Getenv("ACME_HMAC_KEY"))

	// configure TLS via ACME provisioned certificates
	return &tls.Config{
		GetCertificate: (&autocert.Manager{
			Prompt:      autocert.AcceptTOS,
			HostPolicy:  autocert.HostWhitelist(strings.Split(os.Getenv("SERVER_NAMES"), ",")...),
			RenewBefore: 336 * time.Hour, // 14 days

			Client: &acme.Client{
				DirectoryURL: os.Getenv("ACME_DIRECTORY_URL"),
			},

			ExternalAccountBinding: &acme.ExternalAccountBinding{
				KID: os.Getenv("ACME_KID"),
				Key: acmeKey,
			},
		}).GetCertificate,
	}
}
