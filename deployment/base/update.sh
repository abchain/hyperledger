#!/bin/bash

docker build -t abchain/fabric-baseimage:latest update/
docker tag abchain/fabric-baseimage hyperledger/fabric-baseimage:latest
