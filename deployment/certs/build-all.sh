#!/usr/bin/env bash

./build-ca.sh

./build-key.sh gateway
./build-key.sh membersrvc

./copy-certs.sh
