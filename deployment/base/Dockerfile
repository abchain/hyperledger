FROM golang:1.7

RUN apt-get update \
        && apt-get install -y libsnappy-dev zlib1g-dev libbz2-dev \
        && rm -rf /var/cache/apt

# install rocksdb
RUN cd /tmp \
        && git clone --single-branch -b fabric-0.6 --depth 1 https://github.com/noel2004/rocksdb \
        && cd rocksdb \
        && PORTABLE=1 make shared_lib \
        && INSTALL_PATH=/usr/local make install-shared \
        && ldconfig \
        && rm -rf /tmp/rocksdb

# this is only a workaround for current hard-coded problem when using as fabric-baseimage.
RUN ln -s $GOPATH /opt/gopath

# Set workdir and env
WORKDIR $GOPATH

# Copy source files
COPY src/github.com/abchain/fabric $GOPATH/src/github.com/abchain/fabric
COPY src/hyperledger.abchain.org $GOPATH/src/hyperledger.abchain.org

# Build membersrvc, peer 
RUN go install github.com/abchain/fabric/membersrvc
RUN go install github.com/abchain/fabric/peer

# Install membersrvc and peer to $PATH
RUN cp $GOPATH/bin/membersrvc /usr/bin/ \
    && cp $GOPATH/bin/peer /usr/bin/

# Clear
RUN rm -rf $GOPATH/pkg

