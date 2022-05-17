#!/bin/sh
docker build -f Dockerfile.dev -t chat-service-test  . &&
docker rm -f chat-service-test &&
 docker run  --network chat_test_network --entrypoint \
     /opt/app/api/test_entrypoint.sh --name chat-service-test --env-file \
         ./.env.test -v $(pwd):/opt/app/api chat-service-test:latest