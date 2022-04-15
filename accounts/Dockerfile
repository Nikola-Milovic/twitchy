FROM golang:alpine AS build

RUN apk add git

RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go build -o /tmp/accounts ./main.go

FROM alpine:edge

COPY --from=build /tmp/accounts /sbin/accounts

CMD /sbin/accounts