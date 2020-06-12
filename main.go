package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ErrNotImplemented is an error returned when planned functionality is not yet
// implemented.
var ErrNotImplemented = errors.New("not yet implemented")

// ElapsedSince returns the time since a start time in a floating-point number
// of milliseconds.
func ElapsedSince(start time.Time) float32 {
	now := time.Now()
	return float32(now.Sub(start).Nanoseconds()) / float32(time.Millisecond)
}

type ConfigProperties struct {
	logger  *zap.Logger
	zapCfg  *zap.Config
	timeout time.Duration
	mu      sync.Mutex
}

var Config *ConfigProperties

// DefaultConfig creates a new default config
func DefaultConfig() *ConfigProperties {
	config := new(ConfigProperties)

	config.timeout = time.Second * 60

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"stdout",
	}
	cfg.ErrorOutputPaths = []string{
		"stderr",
	}
	cfg.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.MessageKey = "event"
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.StacktraceKey = ""

	config.zapCfg = &cfg
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	config.logger = logger
	return config
}

// Logger returns the configured global logger.
func (c *ConfigProperties) Logger() *zap.Logger {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.logger
}

func (c *ConfigProperties) Timeout() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.timeout
}

func (c *ConfigProperties) SetLevel(level zapcore.Level) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.zapCfg.Level.SetLevel(level)
}

func init() {
	Config = DefaultConfig()
}

func main() {
	// Sync log
	defer Config.logger.Sync()
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "URL")
		return
	}

	t := NewTransport()
	req, _ := http.NewRequest("GET", os.Args[1], nil)

	req = req.WithContext(
		httptrace.WithClientTrace(req.Context(), t.ClientTrace()),
	)

	client, _ := t.Client()
	res, err := client.Do(req)
	if err != nil {
		Config.Logger().Fatal("error",
			zap.Error(err),
		)
	}
	_ = res
	// TODO: do stuff with the response
}
