FROM abchain/fabric-baseimage:latest

# Copy cert files
COPY ca.crt /var/hyperledger/
# Copy config files
COPY membersrvc.yaml $GOPATH/src/github.com/abchain/fabric/membersrvc/
COPY membersrvc.key /var/hyperledger/
COPY membersrvc.crt /var/hyperledger/

