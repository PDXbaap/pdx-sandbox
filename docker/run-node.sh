#!/bin/sh

# Set up docker-in-docker 

/usr/local/bin/wrapdocker

# TODO Start docker controller

# TODO Load docker images 

# Create unprivileged user

pass=$(< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c32)

echo $pass

adduser --no-create-home --disabled-login --gecos '' pdxuser

echo pdxuser:$pass | chpasswd

sudo -u pdxuser /pdx/bin/run-baap.sh
