package main

// Timings is a HTTP Archive (HAR) 1.2 compatible object describing time elapsed
// during various phases of a request-response round trip. All times are
// specified in milliseconds.
//
// Per the HAR spec:
//
// - The Send, Wait and Receive timings are required and must
// have non-negative values.
// - Use -1 for timing values that do not apply the
// request.
//
// See http://www.softwareishard.com/blog/har-12-spec/#timings for more.
type Timings struct {
	Blocked float32 `json:"blocked"`
	DNS     float32 `json:"dns"`
	Connect float32 `json:"connect"`
	Send    float32 `json:"send"`
	Wait    float32 `json:"wait"`
	Receive float32 `json:"receive"`
	SSL     float32 `json:"ssl"`
	Comment string  `json:"comment"`
}

// Trace is a context used for tracking a full round trip of a request,
// inclusive of redirects.
type Trace struct {
	ID      string
	timings []*Timings
}

// Timings returns the timings observed during a Trace.
func (t *Trace) Timings() ([]*Timings, error) {
	return []*Timings{}, ErrNotImplemented
}
