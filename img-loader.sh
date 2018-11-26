#!/bin/bash

#Usage: verifier.sh signer-cert img-dir 

pubk=$(mktemp)

openssl x509 -in $1 -pubkey -noout > $pubk 

for image in $2/*.tgz
do
    sigf="$image.sig"
    if [ -e $sigf ]
    then
	bsig=$(mktemp)
        base64 -d --wrap=0 $sigf >> $bsig
	result=$(openssl dgst -sha256 -verify $pubk -signature $bsig $image)
	rm -rf $bsig
	if [ "$result" == "Verified OK" ] 
	then 
	    docker load < $image
	fi
    fi  
done
