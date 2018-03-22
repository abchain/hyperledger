#!/bin/bash

docker build -t abchain/fabric-baseimage:latest .

docker tag abchain/fabric-baseimage hyperledger/fabric-baseimage:latest
