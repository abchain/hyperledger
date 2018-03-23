#!/usr/bin/env bash

if [ $# != 1 ]; then
    echo "Usage: ./build-key.sh server_name"
    exit
fi

CNF_NAME=fabric-$1.cnf
cp fabric.cnf ${CNF_NAME}
sed -i.bck "s/__ROLE_NAME__/$1/g" ${CNF_NAME}

openssl req -days 3650 -nodes -new -keyout $1.key -out $1.csr -config ${CNF_NAME}
openssl ca -batch -days 3650 -out $1.crt -in $1.csr -config ${CNF_NAME}

rm ${CNF_NAME}
rm ${CNF_NAME}.bck
