#!/usr/bin/env bash

rm index.txt*
touch index.txt
echo 1000 >serial
./build-ca.sh

./build-key.sh gateway
./build-key.sh membersrvc

./copy-certs.sh
