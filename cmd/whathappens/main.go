package main

import (
	"fmt"
	"net/http"
	"net/http/httptrace"
	"os"

	"github.com/travis-g/whathappens"
	"go.uber.org/zap"
)

func main() {
	// Sync log
	defer whathappens.Config.SyncLogger()
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "URL")
		return
	}

	t := whathappens.NewTransport()
	req, _ := http.NewRequest("GET", os.Args[1], nil)

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
}
