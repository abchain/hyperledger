#!/usr/bin/env bash

openssl req -days 3650 -nodes -new -x509 -keyout ca.key -out ca.crt -config fabric.cnf
