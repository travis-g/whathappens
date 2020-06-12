package whathappens

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"

	"go.uber.org/zap"
)

func init() {
	// Reduce log verbosity
	Config.SetLevel(zap.WarnLevel)
}

func Example() {
	req, _ := http.NewRequest("GET", "http://1.1.1.1/", nil)

	// Make a RoundTripper to track the request
	t := NewTransport()

	// Create an HTTP client that uses the transport
	client, _ := t.Client()

	// Set the request to use a tracer and the transport
	req = req.WithContext(
		httptrace.WithClientTrace(req.Context(), t.ClientTrace()),
	)

	// Use the tracer client to make the request. While performing the request,
	// lifecycle hooks will be triggered and timings will be logged
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// The client will have been redirected to the HTTPS site
	fmt.Println(res.Request.URL)
	// Output:
	// https://1.1.1.1/
}
