# PDX sandbox, a setgid docker helper to automate run/stop/stats of PDX sandbox containers 

-- Without setgid

sudo -E ./sandbox docker run -it --rm --memory=100m --cpus="0.1" -v=$PDX_HOME/dapps:/dapps/:ro --name=xzz pdx-sandbox /run-exec.sh test

-- Setgid docker

sudo chgrp docker ./sandbox
sudo chmod g+s ./sandbox

./sandbox docker run -it --rm --memory=100m --cpus="0.1" -v=$PDX_HOME/dapps:/dapps/:ro --name=xzz pdx-sandbox /run-exec.sh test

