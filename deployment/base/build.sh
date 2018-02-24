#!/bin/bash

mkdir -p src
mkdir -p src/github.com/abchain

if [ ! -d src/hyperledger.abchain.org ]; then
    cp "$GOPATH/src/hyperledger.abchain.org" src/ -R
fi

if [ ! -d src/github.com/abchain/fabric ]; then
    cp "$GOPATH/src/github.com/abchain/fabric" src/github.com/abchain -R
fi

docker build -t abchain/fabric-baseimage:latest .

docker tag abchain/fabric-baseimage hyperledger/fabric-baseimage:latest
