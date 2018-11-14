#!/bin/sh

PDX_HOME=/pdx

echo "Starting d-app: " $1

if [ -z "$1" ]
	then	
		echo "No argument supplied, exiting"	
		exit 1

fi

if [ -z "$2" ] 
	then
		echo "Using default security profile" 
		java -jar $PDX_HOME/bin/dapps/java/$1 
	else 
		echo "Using custom security profile" 
		java -Djava.security.manager -Djava.security.policy=$PDX_HOME/bin/dapps/java/$2 -jar $PDX_HOME/bin/dapps/java/$1 
fi

echo "Shutdown d-app: " $1

exit 0

