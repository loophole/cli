# Loophole CLI

https://loophole.cloud

Loophole CLI is one of the available loophole clients.

## Installation

Head over to [the releases page](https://github.com/loophole/cli/releases/latest) and get binary which is suitable for you.

## Quick start

First create an account by executing

```
$ loophole account login
```

and following the instructions there, then execute

```
# Forward application running on local port 3000 to the world
$ loophole http 3000
```

```
# Forward local directory to the world
$ loophole dir ./my-directory
```

Congrats, you can now share the presented link to the world.

For more information head over to [docs](https://loophole.cloud/docs/).


## Development

### Testing

```
$ go test -v ./...
```

### Running

```
# go run main.go
```

### Building

```
# go build -o loophole main.go
```
