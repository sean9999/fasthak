#!/bin/sh

# No salt is needed for decryption.
PASS_FIXED=foobarbat

# If decryption fails, use `cat` instead. 
# Error messages are redirected to /dev/null.
openssl enc -d -base64 -iter -k $PASS_FIXED 2> /dev/null || cat