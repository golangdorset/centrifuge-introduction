# Introduction to Centrifugo
This repository is an introduction to the [Centrifugo][1] real-time messaging
server and [Centrifuge][2] client libraries. This project allows you to build
servers and clients that exchange messages to each other in real time over many
transports such as WebSockets and gRPC.

Centrifugo offers broadcast and subscription based message delivery as well as
ways in which to create secure and private channels.

## Requirements
* [Go 1.13 or newer][3]
* [NodeJS v15 or newer][4]
* Make

## Components
This repository contains 3 components:

* [A Centrifugo server implemented in Go][5]
* [A Centrifuge client implemented in Go][6]
* [A Centrifuge client implemented in NodeJS][7]

Together these components allow you to send messages to all or individual
clients in real-time. Each client gets a user ID which it uses to register and
subscribe to channels on the server.

## Usage

### Build the components
Build all the components by running:

```bash
$ make build
```

Note, the makefile only builds for AMD64 architectures on the following
platforms:

* Linux
* Windows
* Mac (darwin)

It's likely that the code will compile fine for other operating systems and
architectures, but you will need to do this manually with `go build`

### Run the server
Run the server like so:

```bash
$ ./bin/server-[platform]-[architecture]
```

Where `platform` and `architecture` match your current platform. On Windows you
will need to add `.exe`.

The server has the following configuration flags:

```bash
-host string
    IP on which to bind the server
-port string
    Port on which to bind the server (default "8888")
```

### Run the Go client
Run the client like so:

```bash
$ ./bin/client-[platform]-[architecture]
```

Where `platform` and `architecture` match your current platform. On Windows you
will need to add `.exe`.

The client has the following configuration flags:

```bash
-host string
    host of the server (default "localhost")
-port string
    port of the server (default "8888")
-user string
    user ID (default "123")
```

### Run the NodeJS client
The NodeJS client is not compiled ahead of time, so just needs to be run using
this command from the `client/js` directory:

```bash
$ node main.js --host=localhost --port=8888 --user=2
```

The client has the following configuration flags:

```bash
--host string
    host of the server
--port string
    port of the server
--user string
    user ID
```

Unlike the Go client, these flags are mandatory and do not have defaults.

### Send a broadcast message
```bash
$ curl -X POST -i localhost:8888/v1/message/broadcast -H 'Content-Type: application/json' --data-raw '{"foo": "bar"}'
```

You should see this message appear on all clients.

### Send a message to one client
```bash
$ curl -X POST -i 'localhost:8888/v1/message/publish?user=123' -H 'Content-Type: application/json' --data-raw '{"msg": "hello user 123"}'
```

You should see this message appear on just the client registered with user 123

[1]: https://centrifugal.dev/
[2]: https://github.com/centrifugal
[3]: https://go.dev/
[4]: https://nodejs.org/en/
[5]: server
[6]: client/go
[7]: client/js