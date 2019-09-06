#!/usr/bin/env bash


protoc *.proto -I=$GOPATH/src/hyperledger.abchain.org -I=. --go_out=plugins=grpc:.

