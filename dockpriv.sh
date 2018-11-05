#!/bin/sh

echo "start setting privilege for docker: " $1

chgrp docker $1 && chmod g+s $1 

echo "done setting privilege for docker: " $1
