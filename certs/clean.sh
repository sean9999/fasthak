#!/bin/sh

SALT_FIXED=FEEF
PASS_FIXED=foobarbat

openssl enc -base64 -aes-256-ecb -S $SALT_FIXED -k $PASS_FIXED