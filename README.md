# Loophole CLI

https://loophole.cloud

Loophole CLI is one of the available loophole clients.

<a href="https://www.producthunt.com/posts/loophole?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-loophole" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=280773&theme=light" alt="Loophole - Instant hosting, right from your local machine | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>

## Installation

Head over to [the releases page](https://github.com/loophole/cli/releases/latest) and get binary which is suitable for you.

## Quick start

First create an account by executing

```
$ ./loophole account login
```

and following the instructions there, then execute

```
# Forward application running on local port 3000 to the world
$ ./loophole http 3000
```

```
# Forward local directory to the world
$ ./loophole path ./my-directory
```

```
# Forward local directory to the world using WebDAV
$ ./loophole webdav ./my-directory
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
# go run cli.go
```

### Building

```
# go build -o loophole cli.go
```
