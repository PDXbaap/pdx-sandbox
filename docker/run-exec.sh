#!/bin/sh

PDX_HOME=/pdx

if [ -z "$1" ]
	then	
		echo "No exec supplied, exiting"	
		exit 1

fi

echo "Starting d-app: " $1

$PDX_HOME/bin/dapps/exec/$1

echo "Shutdown d-app: " $1

exit 0

