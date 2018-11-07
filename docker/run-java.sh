#!/bin/sh

DAPP_HOME=/dapps

echo "Starting chain code"

if [ -z "$1" ]
	then	
		echo "No argument supplied, exiting"	
		exit 1

fi

if [ -z "$2" ] 
	then
		echo "Using default security profile" 
		java -jar $DAPP_HOME/java/$1 
	else 
		echo "Using custom security profile" 
		java -Djava.security.manager -Djava.security.policy=$DAPP_HOME/jsec/$2 -jar $DAPP_HOME/java/$1 
fi

echo "Shutdown chain code"

exit 0

