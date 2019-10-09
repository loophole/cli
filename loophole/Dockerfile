FROM golang:alpine AS build

WORKDIR /go/src/loophole/

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -o /loophole

#ADD loophole /app/
FROM alpine

WORKDIR /root

COPY ssh ssh
COPY --from=build /loophole /usr/bin/
ENTRYPOINT ["loophole"]