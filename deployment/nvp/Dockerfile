FROM abchain/fabric-baseimage:latest

# Copy cert files
COPY ca.crt /var/hyperledger/
# setup the chaincode sample
COPY core.yaml $GOPATH/src/github.com/abchain/fabric/peer/
COPY gateway.crt /var/hyperledger/
COPY gateway.key /var/hyperledger/

