# dapp-java
#
#       PDX blockchain hypercloud sandbox for JAVA-based smart-contracts
#
#       On PDX blockchain hypercloud, a smart contract has to be sandboxed
#       for isolation and security/resource control, in order to minimize
#       the impact of a potentially harmful smart contract affecting the
#       hosting node and/or other smart contracts on the same host.
#
#       This sandbox applies to openJDK 8 compliant, single-jar smart-contracts
#       executable on barebone alpine Linux systems. 
#
# VERSION 0.1

FROM openjdk:8-jre-alpine

LABEL maintainer="jz@pdx.ltd"

RUN apk add --no-cache bash

COPY run-java.sh /run.sh

ENTRYPOINT ["/run.sh"]

