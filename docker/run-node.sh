#!/bin/sh

# Set up docker-in-docker 

/usr/local/bin/wrapdocker

# Start docker controller

pass=$(< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c32)

echo $pass

adduser --no-create-home --disabled-login --gecos '' pdxuser

echo pdxuser:$pass | chpasswd

sudo -u pdxuser /pdx/bin/run-baap.sh
