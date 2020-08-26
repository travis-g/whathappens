# whathappens

[![GoDoc](https://godoc.org/github.com/travis-g/whathappens?status.svg)][godoc] [![Go Report Card](https://goreportcard.com/badge/github.com/travis-g/whathappens)](https://goreportcard.com/report/github.com/travis-g/whathappens)

A library and CLI for monitoring what happens during an HTTP transaction.

whathappens is a wrapper around the [`httptrace` Golang package][httptrace]. It can be used to instrument health checks, profile network performance and measure SLIs of a service from a client's perspective.

## Phases

```plain
|------------------- Request Duration ---------------------------|
|------------------- Time-to-First-Byte -----------------|
  BLOCKED  DNS             CONNECT             SEND  WAIT RECEIVE
+---------+---+------------------------------+------+----+-------+
|         |   |                              |      |    |       |
|         |   |    +-------------------------+      |    |       |
|         |   |    |   SSL/TLS NEGOTIATION   |      |    |       |
+---------+---+----+-------------------------+------+----+-------+
```

- **Blocked** - Time spent queued for client network resources.
- **DNS** - Time spent resolving a DNS address to an IP.
- **Connect** - Time spent establishing a connection to the server.
- **SSL** - Time spent during the TLS handshake, if the request used SSL/TLS. If included, SSL timings are also included in Connect timings.
- **Send** - Time spent sending the request.
- **Wait** - Time spent waiting for the server to respond.
- **Receive** - Time spent reading the response from the server.

The model aligns with the [HTTP Archive (HAR) 1.2 specification][har-spec] for recording tracing metrics.

<!-- ## Events -->

[godoc]: https://godoc.org/github.com/travis-g/whathappens
[httptrace]: https://blog.golang.org/http-tracing
[har-spec]: http://www.softwareishard.com/blog/har-12-spec/#timings
