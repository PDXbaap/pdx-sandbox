#!/bin/sh

DAPP_HOME=/dapps

if [ -z "$1" ]
	then	
		echo "No exec supplied, exiting"	
		exit 1

fi

echo "Starting chain code: " $1

$DAPP_HOME/exec/$1

echo "Shutdown chain code: " $1

exit 0

