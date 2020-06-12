package whathappens

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

// SyncLogger syncs the underlying logger.
func (c *ConfigProperties) SyncLogger() {
	c.logger.Sync()
}

func init() {
	Config = DefaultConfig()
}
