### Twitchy

Twitchy is a pet project made to play around with new technologies and techniques for writing distributed performant systems. 
The idea of the project is to have a Twitch clone but only have audio instead of video (maybe I'll go for video streaming as well, haven't decided yet)

## Technologies

Languages used: Golang for most of the REST API's, Elixir for real-time chatting
UI: Phoenix LiveView + Tailwind CSS

I've used Fiber, and Chi as web frameworks for the Golang, to try them out. Fiber is really cool and similar to Express, and Chi seems like successor of Gorrila/Mux

Chat microservice has the standard setup, Phoenix + Ecto and no assets/ html.

RabbitMQ as the Message Broker.


## Getting started

Clone the repository, to get it up and running just run the docker compose from the root directory using the helper bash script provided

```
./d-helper.sh up
```

You might need to give permission to the script beforehand. The script is there to merge the compose files for services and databases as it got quite messy considering it's Database-per-Service.


To allow for 🔥blazingly🔥 fast development, I added (Air)[https://github.com/cosmtrek/air] live reload for go apps, and (exsync)[https://github.com/falood/exsync] for the elixir code reloads. 

## Testing

### Chat service
To test the chat service, we have to do it a bit differently than the Golang services. As Elixir likes to spin up services it uses for testing (the Database for example with Ecto), the only solution I could find that didn't slow up the development time was to have a `chat-db-test` and `rabbitmq-test` always running. Then using the bash script `./test.sh` we run the container (mounting our current directory to save some time), the container will run `test_entrypoint.sh` which will wait for `postgresdb` to be ready (not really necessary, could probably remove psql installation from the dockerfile as well and it would work just fine for testing, but I need it for dev, so just leaving in there), run `Credo` to analyze our code, it will reset the database and run the tests. 

Since the dependencies are running beforehand this whole process is quite fast. Also I created a new network just for these tests to make sure they don't mess up anything outside of the tests. 

```
docker network create chat-test-network // create the network used in the tests to connected the services
./start_test_deps.sh  // start the services in -d mode and have them keep running in the background
./test.sh // clean/build and run our container used to test the code
```

### Other services
Other services are written in Golang and are purely unit tests, to run them all we have to do is run the following command in their directory
```
go test ./... 
```

## Notes

Code quality is mediocre, as I was trying out a bunch of things I sometimes would lazy out on some aspects of the architecture or "good practices". It could use a bit of refactoring here and there.

The services are missing a few layers, most notably a repository or a data layer. It would be much cleaner and easier to change databases later, but for my purposes I am keeping it easier to change stuff for the time being.

There is a k8 folder, I played around with Kubernetes and Skaffold to get a feel for them, but didn't want to probe any deeper at the moment as my focus is purely on architecture/ actual Microservices.

### Improvements
- [ ] Workbox pattern, store events and the action that emits them in the same transaction to make sure they both succeed
- [ ] Add better pooling for workers
- [ ] Repository layer/ data layer, clearer separation of layers in the project
- [ ] CQRS or event sourcing could be added at a later stage, will have to play around with that when the need arises, currently it's too simple of a project
- [ ] Rename files and move interfaces around, "I" interface naming is unconventional (except in Java...). Eg model -> domain, client isn't the best name for the RabbitMQ Consumer either