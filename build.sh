#!/bin/bash

SEMVER="$(git tag --sort=-version:refname | head -n 1)"

#	@note: certs should be renewed every so often: https://www.rec.la/

#	build the binary
go build -v \
	-ldflags="-X 'main.Version=$SEMVER' -X 'app/build.User=$(id -u -n)' -X 'app/build.Time=$(date)'" \
	-o ./build/
