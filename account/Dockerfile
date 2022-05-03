FROM golang:alpine AS build

RUN apk add git

RUN mkdir /src
RUN mkdir /common_go
ADD ./account /src
ADD ./common_go /common_go
WORKDIR /src

RUN go build -o /tmp/account ./main.go

FROM alpine:edge

COPY --from=build /tmp/account /sbin/account

RUN mkdir -p /sbin/db/migrations

COPY --from=build /src/db/migrations /sbin/db/migrations

EXPOSE $PORT

CMD /sbin/account