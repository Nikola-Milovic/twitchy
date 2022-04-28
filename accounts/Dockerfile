FROM golang:alpine AS build

RUN apk add git

RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go build -o /tmp/accounts ./main.go

FROM alpine:edge

COPY --from=build /tmp/accounts /sbin/accounts

RUN mkdir -p /sbin/db/migrations

COPY --from=build /src/db/migrations /sbin/db/migrations

EXPOSE $PORT

CMD /sbin/accounts