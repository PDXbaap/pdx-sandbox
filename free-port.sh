#!/bin/sh

dirname "$(readlink -f "$0")"

echo $DIR

port=$1

isfree=$(netstat -t4ln | grep LISTEN | grep ":$port ")


while [ -n "$isfree" ]; do 
	port=$(( $port + 1 ))
	isfree=$(netstat -t4ln | grep LISTEN | grep ":$port ")
done

echo "usable port: $port"


