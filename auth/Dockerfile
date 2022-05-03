FROM golang:alpine AS build

RUN apk add git

RUN mkdir /src
RUN mkdir /common_go
ADD ./auth /src
ADD ./common_go /common_go
WORKDIR /src

RUN go build -o /tmp/auth ./main.go

FROM alpine:edge

COPY --from=build /tmp/auth /sbin/auth

RUN mkdir -p /sbin/db/migrations

COPY --from=build /src/db/migrations /sbin/db/migrations

EXPOSE $PORT

CMD /sbin/auth