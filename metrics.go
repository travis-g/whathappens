package whathappens

import (
	"github.com/armon/go-metrics"
)

var (
	// Sink is the package-wide metrics sink.
	Sink *metrics.Metrics
)

func init() {
	var err error
	sink := metrics.NewInmemSink(Config.metricsInterval, Config.retainInterval)
	// TODO: allow use of a custom prefix/metrics config
	Sink, err = metrics.New(metrics.DefaultConfig("whathappens"), sink)
	if err != nil {
		panic(err)
	}
}
