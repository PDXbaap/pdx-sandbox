#!/bin/bash

#Usage: verifier.sh signer-cert file signature  

pubk=$(mktemp)
sigf=$(mktemp)

openssl x509 -in $1 -pubkey -noout > $pubk 

if [ -e $3 ] 
then 
    base64 -d --wrap=0 $3 >> $sigf
else 
    echo $3 | base64 -d --wrap=0 > $sigf 
fi

result=$(openssl dgst -sha256 -verify $pubk -signature $sigf $2)

rm -rf $pubk $sigf

echo $result
