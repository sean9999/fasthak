#!/bin/bash

SEMVER="$(git tag --sort=-version:refname | head -n 1)"

#	create the certificates. They will be embedded in the binary
#	should create localhost.pem and locahost-key.pem
#	the build will fail if those files don't exist.
#	@see: https://github.com/FiloSottile/mkcert#readme

#	mkcert -ecdsa localhost

#	build the binary
go build -v \
	-ldflags="-X 'main.Version=$SEMVER' -X 'app/build.User=$(id -u -n)' -X 'app/build.Time=$(date)'" \
	-o ./build/
