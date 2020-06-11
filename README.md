# whathappens

A library and CLI for figuring out what happens during an HTTP transaction.

whathappens is a wrapper around the `httptrace` Golang library. It can be used to instrument health checks, profile network performance and measure SLIs of a service from a client's perspective.

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
