#!/bin/bash

#	create the certificates. They will be embedded in the binary
#	should create localhost.pem and locahost-key.pem
#	@see: https://github.com/FiloSottile/mkcert#readme

#	mkcert -ecdsa localhost

#	build the binary
go build -o ./build/
