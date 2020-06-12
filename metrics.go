package main

import (
	"github.com/armon/go-metrics"
)

var (
	// Sink is the package-wide metrics sink.
	Sink *metrics.Metrics
)

func init() {
	sink, err := metrics.NewStatsdSink("localhost:8125")
	if err != nil {
		panic(err)
	}
	// TODO: allow use of a custom prefix/metrics config
	Sink, err = metrics.New(metrics.DefaultConfig("whathappens"), sink)
	if err != nil {
		panic(err)
	}
}
