#!/bin/bash
IS_RUNNING=`docker-compose ps -q`
if [[ "$IS_RUNNING" != "" ]]; then
    docker-compose down
fi

docker-compose up -d

# This commands should be run only for initialization only


echo "setup the replication sets"
docker-compose exec mongo-db1 sh -c "mongo --port 27017 < scripts/init-rs.js

echo "wait 20 seconds until replication sets are ready"
sleep 20

echo "create user"
docker-compose exec mongo-db1 sh -c "mongo --port 27017 < scripts/create_user.js"