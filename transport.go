package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
)

// Check implements at compile time
var _ http.RoundTripper = (*Transport)(nil)

// Transport is an http.RoundTripper that keeps track of the in-flight request
// and implements hooks to report HTTP tracing events. A Transport should only
// be used once per HTTP transaction but reused if following redirects.
//
// A Transport's ID and the IDs of requests made using it are RFC 4122 UUIDs.
//
// FIXME: Transports following redirects destroy previous request data.
type Transport struct {
	coreID           string
	transportID      string
	current          *http.Request
	currentRequestID string
	startTime        time.Time
	connStartTime    time.Time
	connectStartTime time.Time
	dnsStartTime     time.Time
	tlsStartTime     time.Time

	timings []*Timings
	labels  []metrics.Label

	mu sync.RWMutex
}

type TransportLog struct {
	Time      time.Time
	Transport string
	Request   string
	Hook      string
}

// ErrNilTransport is returned when a nil Transport is referenced.
var ErrNilTransport = errors.New("nil transport")

// NewTransport allocates a transport and assigns it a UUID.
func NewTransport() (t *Transport) {
	t = new(Transport)
	coreBytes, err := uuid.GenerateRandomBytes(10)
	if err != nil {
		panic(err)
	}
	t.coreID = fmt.Sprintf("%x-%x-%x-%x",
		coreBytes[0:4],
		coreBytes[4:6],
		coreBytes[6:8],
		coreBytes[8:10],
	)
	transportBytes, err := uuid.GenerateRandomBytes(6)
	if err != nil {
		panic(err)
	}
	t.transportID = fmt.Sprintf("%s-%x",
		t.coreID,
		transportBytes,
	)
	return
}

// Timings returns the timings observed by the Transport as a slice of Timings.
func (t *Transport) Timings() ([]*Timings, error) {
	if t == nil {
		return []*Timings{}, ErrNilTransport
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return []*Timings{}, ErrNotImplemented
}

// RoundTrip implements http.RoundTripper and wraps
// http.DefaultTransport.RoundTrip to keep track of the current request.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t == nil {
		return nil, ErrNilTransport
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	// build request UUID
	currentRequestIDBytes, err := uuid.GenerateRandomBytes(6)
	if err != nil {
		panic(err)
	}
	t.currentRequestID = fmt.Sprintf("%s-%x", t.coreID, currentRequestIDBytes)
	Config.Logger().Info("RoundTripStart",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.String("url", req.URL.String()),
		zap.String("method", req.Method),
		zap.String("thing", req.RemoteAddr),
	)
	t.current = req
	t.startTime = time.Now()
	return http.DefaultTransport.RoundTrip(req)
}

// ClientTrace returns an httptrace.ClientTrace that performs the given timings
// when its hooks are triggered.
func (t *Transport) ClientTrace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn:           t.GetConn,
		DNSStart:          t.DNSStart,
		DNSDone:           t.DNSDone,
		ConnectStart:      t.ConnectStart,
		TLSHandshakeStart: t.TLSHandshakeStart,
		TLSHandshakeDone:  t.TLSHandshakeDone,
		ConnectDone:       t.ConnectDone,
		GotConn:           t.GotConn,

		// TODO
		WroteHeaderField: nil,
		WroteHeaders:     nil,
		Wait100Continue:  nil,
		WroteRequest:     nil,

		Got1xxResponse:       t.Got1xxResponse,
		GotFirstResponseByte: t.GotFirstResponseByte,
	}
}

// Client returns an http.Client that will use the transport.
func (t *Transport) Client() (*http.Client, error) {
	if t == nil {
		return nil, ErrNilTransport
	}
	return &http.Client{
		Transport: t,
		Timeout:   Config.Timeout(),
		// TODO: CheckRedirect
	}, nil
}

func (t *Transport) DNSStart(info httptrace.DNSStartInfo) {
	now := time.Now()
	Config.Logger().Info("DNSStart",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.Time("time", now),
		zap.Any("info", info),
	)
	t.dnsStartTime = now
}

func (t *Transport) DNSDone(info httptrace.DNSDoneInfo) {
	duration := ElapsedSince(t.dnsStartTime)
	Config.Logger().Info("DNSDone",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.Float32("duration", duration),
		zap.Any("info", info),
	)
	Sink.AddSample([]string{"dns"}, duration)
}

func (t *Transport) TLSHandshakeStart() {
	now := time.Now()
	t.tlsStartTime = now
	Config.Logger().Info("TLSHandshakeStart",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
	)
}

func (t *Transport) TLSHandshakeDone(state tls.ConnectionState, err error) {
	duration := ElapsedSince(t.tlsStartTime)
	defer Sink.AddSample([]string{"tls", "handshake"}, duration)
	Config.Logger().Info("TLSHandshakeDone",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.Float32("duration", duration),
		zap.Error(err),
	)
}

func (t *Transport) GetConn(hostport string) {
	t.connStartTime = time.Now()
	Config.Logger().Info("GetConn",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.String("addr", hostport),
	)
}

func (t *Transport) GotConn(info httptrace.GotConnInfo) {
	duration := ElapsedSince(t.connStartTime)
	defer Sink.AddSample([]string{"connect", "open"}, duration)
	Config.Logger().Info("GotConn",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.Float32("duration", duration),
		zap.String("addr", info.Conn.RemoteAddr().String()),
		zap.Bool("reused", info.Reused),
	)
}

func (t *Transport) ConnectStart(network, addr string) {
	now := time.Now()
	t.connectStartTime = now
	Config.Logger().Info("ConnectStart",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.String("network", network),
		zap.String("addr", addr),
	)
}

func (t *Transport) ConnectDone(network, addr string, err error) {
	duration := ElapsedSince(t.connectStartTime)
	defer Sink.AddSample([]string{"connect", network}, duration)
	Config.Logger().Info("ConnectDone",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.Float32("duration", duration),
		zap.String("network", network),
		zap.String("addr", addr),
		zap.Error(err),
	)
}

func (t *Transport) GotFirstResponseByte() {
	elapsed := ElapsedSince(t.startTime)
	defer Sink.AddSample([]string{"time_to_first_byte"}, elapsed)
	Config.Logger().Info("GetFirstResponseByte",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.Float32("elapsed", elapsed),
	)
}

func (t *Transport) Got1xxResponse(code int, _ textproto.MIMEHeader) error {
	elapsed := ElapsedSince(t.startTime)
	defer Sink.AddSample([]string{"1xx_response"}, elapsed)
	Config.Logger().Info("Got1xxResponse",
		zap.String("transport", t.transportID),
		zap.String("request", t.currentRequestID),
		zap.Float32("elapsed", elapsed),
		zap.Int("code", code),
	)
	return nil
}
