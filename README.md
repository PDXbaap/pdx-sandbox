# PDX sandboxer, a setgid docker helper to automate run/stop/stats of PDX sandbox containers 

-- Without setgid

sudo -E ./sandboxer docker run -it --rm --memory=100m --cpus="0.1" -v=$PDX_HOME/dapps:/dapps/:ro --name=xzz pdx-dapp-omni /run-exec.sh test

-- Setgid docker

sudo chgrp docker ./sandboxer
sudo chmod g+s ./sandboxer

./sandboxer docker run -it --rm --memory=100m --cpus="0.1" -v=$PDX_HOME/dapps:/dapps/:ro --name=xzz pdx-dapp-omni /run-exec.sh test

