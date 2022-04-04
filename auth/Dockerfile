FROM golang:alpine AS build

RUN apk add git

RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go build -o /tmp/auth ./main.go

FROM alpine:edge

COPY --from=build /tmp/auth /sbin/auth

CMD /sbin/auth