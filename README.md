# PDX sandboxer, a setuid docker helper to start/stop PDX docker containers

-- Without setuid/gid

sudo -E ./sandboxer docker run -it --rm --memory=100m --cpus="0.1" -v $PDX_HOME/dapps:/dapps/ --name xzz pdx-dapp-omni /run-exec.sh test

-- Setgid docker

sudo chgrp docker ./sandboxer
sudo chmod g+s ./sandboxer

./sandboxer docker run -it --rm --memory=100m --cpus="0.1" -v $PDX_HOME/dapps:/dapps/ --name xzz pdx-dapp-omni /run-exec.sh test

