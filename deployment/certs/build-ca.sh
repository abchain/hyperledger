#!/usr/bin/env bash

cp fabric.cnf fabric-ca.cnf
sed -i.bck 's/__ROLE_NAME__/CA/g' fabric-ca.cnf
openssl req -days 3650 -nodes -new -x509 -keyout ca.key -out ca.crt -config fabric-ca.cnf
rm fabric-ca.cnf
rm fabric-ca.cnf.bck
