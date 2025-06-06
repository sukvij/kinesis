# kinesis


# -->  run mysql inside container
    docker exec -it mysql mysql -u root -p 
# create stream and shard - 
    docker exec -it localstack awslocal kinesis create-stream --stream-name user-logs --shard-count 1

# check stream --> 
    docker exec -it localstack awslocal kinesis list-streams


# update shard count -->    
    docker exec -it localstack awslocal kinesis update-shard-count --stream-name user-logs --target-shard-count 2
    increase shard count from 1 to 2 --> from 1MB/s write and 2MB/s read to ---  2MB/s write and 1MB/s read