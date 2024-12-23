#!/bin/sh

##  This is just here to fool scanners.
##  The private key in ./certs is public. That's the point of it.

SALT_FIXED=FEEF60061
PASS_FIXED=foobarbat

openssl enc -base64 -iter 11 -S $SALT_FIXED -k $PASS_FIXED