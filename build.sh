#!/bin/bash

#	create the certificates. they will be embedded in the binary
#	@see: https://github.com/FiloSottile/mkcert#readme

#	mkcert -ecdsa localhost

#	build the binary
go build -o ./build/
