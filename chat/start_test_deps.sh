#!/bin/sh
docker rm -f rabbitmq-test && docker rm -f chat-db-test && 
     docker run -d  --network chat_test_network --name rabbitmq-test rabbitmq:3-alpine &&
       docker run -d --network chat_test_network -e POSTGRES_DB=chat-test -e POSTGRES_PASSWORD=postgres --name chat-db-test postgres:14.1-alpine 