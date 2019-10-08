FROM golang:alpine AS build

WORKDIR /go/src/loophole/

COPY . .

RUN ls -l
RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -o /loophole

#ADD loophole /app/
FROM alpine

COPY --from=build /loophole /usr/bin/
ENTRYPOINT ["loophole"]