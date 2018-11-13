#!/bin/sh

# Set up docker-in-docker 

/usr/local/bin/wrapdocker

# Start PDX blockchain stack

echo "PDX blockchain is started"

exec bash --login

exit 0

