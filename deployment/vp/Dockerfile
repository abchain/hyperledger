FROM abchain/fabric-baseimage:latest

# Copy cert files
COPY ca.crt /var/hyperledger/
# Copy config files
COPY core.yaml $GOPATH/src/github.com/abchain/fabric/peer/
COPY noops_config.yaml $GOPATH/src/github.com/abchain/fabric/consensus/nooops/config.yaml
COPY pbft_config.yaml $GOPATH/src/github.com/abchain/fabric/consensus/pbft/config.yaml

