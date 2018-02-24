#!/usr/bin/env bash

if [ $# != 1 ]; then
    echo "Usage: ./build-key.sh server_name"
    exit
fi

openssl req -days 3650 -nodes -new -keyout $1.key -out $1.csr -config fabric.cnf

openssl ca -days 3650 -out $1.crt -in $1.csr -config fabric.cnf
