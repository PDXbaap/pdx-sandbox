#!/bin/sh

# Start PDX blockchain stack
#
# Note that all docker cmds via the privileged docker controller 
#
# 1) Start PDX BaaP
# 2) If first time, initialize PDX blockchain hypercloud "node"
# 	2.1) Create key store @ $PDX_HOME/data/node/secrets ($PDX_NODE_INFO)
#	     ./pdxc --datadir $PDX_NODE_INFO account new
# 	2.2) Create node key
#	     ./bootnode -genkey $PDX_NODE_INFO nodekey
#	2.3) Generate enode URL to be used for node registration
#	     ./bootnode -nodekey $PDX_NODE_INFO/nodekey -writeaddress
#       2.4) Register node to centralized IaaS service via REST (node key, miner public key, IP, etc)
#	     IaaS return the trust chain to join with the following info:
#		id, chain-type, engine-type, genesis, static nodes, ...
#	2.5) Prepare trust chain data directory ($PDX_CHAIN_DATA)
#		create $PDX_HOME/data/chains/{chain-id}
#		copy genesis and static nodes
#	2.6) Check if chainstack docker image is available.
#	     Convention {engine-type}-chainstack:latest, e.g. pdx-chainstack:latest
#	     If not, docker restore it from $PDX_HOME/bin/images/{engine-type}-chainstack.tgz
#	2.8) Select rpc-port, p2p-port and grpc-port
#       2.7) Start chain docker container
#	     docker run -v=$PDX_NODE_INFO:/pdx/data/node/secrets:ro -v=$PDX_CHAIN_DATA):/pdx/data/chain:rw \
#		--env rpc_port={rpc-port} --env p2p_port={p2p-port} --env grpc_port={grpc-port} --name pdx-chain-{chainid} {engine-type}-chainstack:latest
#	2.8) Call PDX BaaP REST API to notify new chain (chain-id, rpc-port, grpc-port)
#	2.9) Call centralized IaaS service to update chain-node binding
#	2.10) Update $PDX_HOME/data/node/chaininfo.conf
#
# 3) If not first time, for each chain from $PDX_HOME/data/node/chaininfo.conf
#       3.1) Start chain docker container
#	     docker run -v=$PDX_NODE_INFO:/pdx/data/node/secrets:ro -v=$PDX_CHAIN_DATA):/pdx/data/chain:rw \
#		--env rpc_port={rpc-port} --env p2p_port={p2p-port} --env grpc_port={grpc-port} --name pdx-chain-{chainid} {engine-type}-chainstack:latest
#	3.2) Call PDX BaaP REST API to notify new chain (chain-id, rpc-port, grpc-port)
#	3.3) Call centralized IaaS service to update chain-node binding
#

echo "PDX blockchain is started"

exec /bin/bash --login

exit 0

