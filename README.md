# whathappens

[![GoDoc](https://godoc.org/github.com/travis-g/whathappens?status.svg)][godoc] [![Go Report Card](https://goreportcard.com/badge/github.com/travis-g/whathappens)](https://goreportcard.com/report/github.com/travis-g/whathappens)

A library and CLI for monitoring what happens during an HTTP transaction.

whathappens is a wrapper around the [`httptrace` Golang package][httptrace]. It can be used to instrument health checks, profile network performance and measure SLIs of a service from a client's perspective.

## Request-Response Phases

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

[godoc]: https://godoc.org/github.com/travis-g/whathappens
[httptrace]: https://blog.golang.org/http-tracing
