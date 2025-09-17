package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-clean-arch/internal/shared/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.Logger {
	config := zap.NewProductionConfig()

	// Set log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, _ := config.Build()
	return logger
}

func NewWithLoki(level string, lokiConfig config.LokiConfig) *zap.Logger {
	// Create base logger
	baseLogger := New(level)

	if !lokiConfig.Enabled {
		return baseLogger
	}

	// Create custom core that writes to both console and Loki
	lokiCore := &LokiCore{
		lokiURL: lokiConfig.URL,
		labels:  lokiConfig.Labels,
		enc:     zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
	}

	// Combine cores
	core := zapcore.NewTee(
		baseLogger.Core(),
		lokiCore,
	)

	return zap.New(core)
}

type LokiCore struct {
	lokiURL string
	labels  map[string]string
	enc     zapcore.Encoder
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LokiRequest struct {
	Streams []LokiStream `json:"streams"`
}

func (c *LokiCore) Enabled(level zapcore.Level) bool {
	return true
}

func (c *LokiCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	for _, field := range fields {
		clone.enc.AddString(field.Key, fmt.Sprintf("%v", field.Interface))
	}
	return clone
}

func (c *LokiCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c *LokiCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}

	// Send to Loki asynchronously
	go func() {
		defer buf.Free()
		c.sendToLoki(entry.Time, buf.String())
	}()

	return nil
}

func (c *LokiCore) sendToLoki(timestamp time.Time, message string) {
	// Create Loki request
	lokiReq := LokiRequest{
		Streams: []LokiStream{
			{
				Stream: c.labels,
				Values: [][]string{
					{
						fmt.Sprintf("%d", timestamp.UnixNano()),
						message,
					},
				},
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(lokiReq)
	if err != nil {
		return
	}

	// Send HTTP request to Loki
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", c.lokiURL+"/loki/api/v1/push", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client.Do(req)
}

func (c *LokiCore) Sync() error {
	return nil
}

func (c *LokiCore) clone() *LokiCore {
	return &LokiCore{
		lokiURL: c.lokiURL,
		labels:  c.labels,
		enc:     c.enc.Clone(),
	}
}
