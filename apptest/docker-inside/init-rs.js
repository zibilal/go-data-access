rs.initiate({
    "_id": "myreplicaset",
    "members": [
        {
            "_id": 0,
            "host": "mongo-db1:27017"
        },
        {
            "_id": 1,
            "host": "mongo-db2:27017"
        },
        {
            "_id": 2,
            "host": "mongo-db3:27017"
        }
    ]
})