#!/bin/sh

SALT_FIXED=FEEF60061
PASS_FIXED=foobarbat

openssl enc -base64 -iter -S $SALT_FIXED -k $PASS_FIXED