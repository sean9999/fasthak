package certs

import (
	_ "embed"
)

//go:embed backloop.dev-ca.crt
var Authority []byte

//go:embed backloop.dev-cert.crt
var Cert []byte

//go:embed backloop.dev-key.pem
var Key []byte
