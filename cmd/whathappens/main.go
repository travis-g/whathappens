package main

import (
	"net/http"
	"net/http/httptrace"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/travis-g/whathappens"
	"go.uber.org/zap"
)

var (
	wait  int
	count int
	// level zapcore.Level
)

func makeRequest(url string) error {
	t := whathappens.NewTransport()
	req, _ := http.NewRequest("GET", url, nil)

	req = req.WithContext(
		httptrace.WithClientTrace(req.Context(), t.ClientTrace()),
	)

	client, _ := t.Client()

	res, err := client.Do(req)
	if err != nil {
		whathappens.Config.Logger().Fatal("error",
			zap.Error(err),
		)
	}
	_ = res
	// TODO: do stuff with the response
	return err
}

func main() {
	// Sync log
	defer whathappens.Config.SyncLogger()

	flag.IntVarP(&wait, "wait", "i", 5, "wait")
	flag.IntVarP(&count, "count", "c", 1, "count")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}

	requests := 0
	ticker := time.NewTicker(time.Duration(wait) * time.Second)

	for ; true; <-ticker.C {
		_ = makeRequest(args[0])
		if count > 0 {
			requests++
			if requests >= count {
				break
			}
		}
	}
}
