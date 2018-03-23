FROM abchain/fabric-baseimage:latest

RUN go get -u hyperledger.abchain.org/utils
RUN go get -u github.com/abchain/fabric/peer
RUN go install github.com/abchain/fabric/membersrvc

# Install membersrvc and peer to $PATH
RUN cp $GOPATH/bin/membersrvc /usr/bin/ \
    && cp $GOPATH/bin/peer /usr/bin/

# Clear
RUN rm -rf $GOPATH/pkg

