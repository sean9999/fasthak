#!/bin/bash

##  @description: run the server in dev mode.
##	production mode would look like this:
##	fasthak -dir=public --port=9443

go run *.go --dir=public
