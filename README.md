# PDX sandbox, a secure privileged service for hardcore docker sandboxing.

PDX sandbox supports signed docker images and fine-grained access control. 

Supported docker commands: run, load, stop & stat

## Dependencies

- Linux 4.15 x64 
- docker-ce
- openssl
- PDX iaas-compute installed on $PDX_HOME

## Configuration

PDX sandbox assumes $PDX_HOME, installation root of PDX iaas-compute has the following directory structure:

├── bin
│   ├── ........... 
│   ├── dapp
│   │   ├── exec
│   │   │   ├── test-1.0.0
│   │   │   ├── test-1.0.0.pub
│   │   │   └── test-1.0.0.sig
│   │   └── java
│   │       ├── HelloWorld-1.0.0.jar
│   │       ├── HelloWorld-1.0.0.jar.pub
│   │       ├── HelloWorld-1.0.0.jar.sig
│   │       └── java-default.policy
│   ├── image
│   │   ├── pdx-appsandbox-1.0.0.tgz
│   │   ├── pdx-appsandbox-1.0.0.tgz.sig
│   │   ├── pdx-chainstack-1.0.0.tgz
│   │   └── pdx-chainstack-1.0.0.tgz.sig
│   ├── img-loader.sh
│   ├── img-signer.sh
│   ├── img-verify.sh
│   └── sandbox
├── chain
│   └── ........... 
├── conf
│   └── signer.crt
├── dapp
│   └── ...........
├── node
│   └── ........... 
└── temp
    ├── sandbox.data
    └── sandbox.lock

## Starting sandbox

./sandbox -addr=127.0.0.1:7391 -home=$PDX_HOME

## Testing using curl

  Do NOT use -it, -i or it.

  curl -d "docker run --rm --memory=100m --cpus=0.1 -v=$PDX_HOME/dapps:/dapps/:ro --name=xzzz pdx-sandbox /bin/sh" -X POST http://localhost:7391

  curl -d "docker load --input pdx-appsandbox-1.0.0.tgz" -X POST http://localhost:7391

## Certificate & verification

  openssl ecparam -genkey -name secp256k1 -out signer.key
  openssl req -new -sha256 -key signer.key -out signer.csr

	Country Name (2 letter code) [AU]:US
	State or Province Name (full name) [Some-State]:California
	Locality Name (eg, city) []:San Jose
	Organization Name (eg, company) [Internet Widgits Pty Ltd]:PDX Technologies, Inc.
	Organizational Unit Name (eg, section) []:Blockchain Hypercloud
	Common Name (e.g. server FQDN or YOUR name) []:signer.pdx.link
	Email Address []:jz@pdx.ltd

  openssl req -x509 -sha256 -days 3650 -key signer.key -in signer.csr -out signer.crt

- Sign & verify, see the img-signer/verify.sh for details

  ./img-signer.sh ./signer.key /home/jz/pdx-home/bin/image/pdx-chainstack-1.0.0.tgz /home/jz/pdx-home/bin/images/pdx-chainstack-1.0.0.tgz.sig

  ./img-verify.sh ./signer.crt /home/jz/pdx-home/bin/image/pdx-appsandbox-1.0.0.tgz /home/jz/pdx-home/bin/images/pdx-appsandbox-1.0.0.tgz.sig
Verified OK

  ./img-loader.sh ./signer.crt $PDX_HOME/bin/image
