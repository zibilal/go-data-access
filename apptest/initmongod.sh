#!/bin/bash

[ ! -f ~/mongodb/path/db1 ] && mkdir -p ~/mongodb/path/db1
[ ! -f ~/mongodb/path/db2 ] && mkdir -p ~/mongodb/path/db2
[ ! -f ~/mongodb/path/db3 ] && mkdir -p ~/mongodb/path/db3

mongod --replSet "rs0" --port 27027 --dbpath ~/mongodb/path/db1 --bind_ip 127.0.0.1 &
mongod --replSet "rs0" --port 27037 --dbpath ~/mongodb/path/db2 --bind_ip 127.0.0.1 &
mongod --replSet "rs0" --port 27047 --dbpath ~/mongodb/path/db3 --bind_ip 127.0.0.1 &

echo "Sleep for twenty seconds"
sleep 30

mongo --port 27027 < init-rs-by-port.js