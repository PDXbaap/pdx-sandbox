#!/bin/bash

#Usage: signer.sh signer-key file-to-sign sig-file 
#	
openssl dgst -sha256 -sign $1 $2 | base64 --wrap=0 > $3 

