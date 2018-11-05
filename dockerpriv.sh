#!/bin/sh

echo ""

echo "PDX blockchain hypercloud uses docker to sandbox smart-contracts"
echo "for isolation and security/resource controls."

echo ""

echo "However, docker requires root or run-as-docker-group priviledge"
echo "to function."

echo ""

echo "Here we grant docker privilege to PDX sandboxer so that on-demand"
echo "automatic sandboxing can be performed."

echo ""

echo "Access control against abuse is done by the open-sourced PDX"
echo "sandboxer, which starts or stops docker instances on-demand."

echo ""

echo "Please visit https://github.com/PDXbaap/pdx-sandboxer for more info."


echo ""
echo ""

if [ -z "$1" ]
	then
		echo "path to pdx sandboxer not supplied, exiting\n"
		exit 1

fi

echo "start setting privilege for docker access: " $1

chgrp docker $1 && chmod g+s $1 

echo "done setting privilege for docker access: " $1
