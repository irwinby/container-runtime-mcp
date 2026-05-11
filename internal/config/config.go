package config

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/irwinby/container-runtime-mcp/pkg/logger"
	env "github.com/sethvargo/go-envconfig"
)

const prefix = "CONTAINER_RUNTIME_"

type TransportType string

const (
	TransportStdio TransportType = "stdio"
	TransportHTTP  TransportType = "http"
)

type TransportConfig struct {
	Type TransportType `env:"MCP_TRANSPORT,default=stdio"`
	HTTP *HTTPTransportConfig
}

type HTTPTransportConfig struct {
	Addr           string        `env:"MCP_HTTP_ADDR,default=127.0.0.1:8080"`
	Path           string        `env:"MCP_HTTP_PATH,default=/mcp"`
	SessionTimeout time.Duration `env:"MCP_HTTP_SESSION_TIMEOUT,default=30m"`
	ReadTimeout    time.Duration `env:"MCP_HTTP_READ_TIMEOUT,default=10s"`
	IDLETimeout    time.Duration `env:"MCP_HTTP_IDLE_TIMEOUT,default=120s"`
	AuthToken      string        `env:"MCP_HTTP_AUTH_TOKEN,default="`
}

type MCPServer struct {
	Name            string `env:"MCP_SERVER_NAME,default=Container Runtime"`
	Title           string `env:"MCP_SERVER_TITLE,default="`
	Version         string `env:"MCP_SERVER_VERSION,default=1.0.0"`
	TransportConfig TransportConfig
}

type TelemetryConfig struct {
	Enabled      bool          `env:"MCP_TELEMETRY_ENABLED,default=false"`
	Addr         string        `env:"MCP_TELEMETRY_ADDR,default=127.0.0.1:9090"`
	PPROFEnabled bool          `env:"MCP_TELEMETRY_PPROF_ENABLED,default=false"`
	ReadTimeout  time.Duration `env:"MCP_TELEMETRY_READ_TIMEOUT,default=10s"`
	IDLETimeout  time.Duration `env:"MCP_TELEMETRY_IDLE_TIMEOUT,default=120s"`
}

type Config struct {
	MCPServer
	RemoteOperationTimeout time.Duration `env:"MCP_REMOTE_OPERATION_TIMEOUT,default=10m"`
	ReadOnly               bool          `env:"MCP_READ_ONLY,default=false"`
	LogLevel               logger.Level  `env:"LOG_LEVEL,default=info"`
	Telemetry              TelemetryConfig
}

// Validate checks the configuration for errors.
func (cfg *Config) Validate() error {
	if cfg.RemoteOperationTimeout < 0 {
		return fmt.Errorf("remote operation timeout must not be negative")
	}

	if err := cfg.validateLogLevel(); err != nil {
		return err
	}

	if err := cfg.validateTransport(); err != nil {
		return err
	}

	if err := cfg.validateTelemetry(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) validateLogLevel() error {
	switch cfg.LogLevel {
	case logger.DebugLevel, logger.InfoLevel, logger.WarnLevel, logger.ErrorLevel:
		return nil
	default:
		return fmt.Errorf("unsupported log level: %s", cfg.LogLevel)
	}
}

func (cfg *Config) validateTransport() error {
	switch cfg.TransportConfig.Type {
	case TransportStdio:
		return nil
	case TransportHTTP:
		return cfg.validateHTTP()
	default:
		return fmt.Errorf("unsupported transport: %s", cfg.TransportConfig.Type)
	}
}

func (cfg *Config) validateHTTP() error {
	if cfg.TransportConfig.HTTP == nil {
		return fmt.Errorf("http transport requires http configuration")
	}

	if !strings.HasPrefix(cfg.TransportConfig.HTTP.Path, "/") {
		return fmt.Errorf("http path must start with /")
	}

	if cfg.TransportConfig.HTTP.SessionTimeout < 0 {
		return fmt.Errorf("http session timeout must not be negative")
	}

	if cfg.TransportConfig.HTTP.ReadTimeout < 0 {
		return fmt.Errorf("http read timeout must not be negative")
	}

	if cfg.TransportConfig.HTTP.IDLETimeout < 0 {
		return fmt.Errorf("http idle timeout must not be negative")
	}

	if !cfg.ReadOnly && cfg.TransportConfig.HTTP.AuthToken == "" {
		isLocal, err := isLocalHTTPAddr(cfg.TransportConfig.HTTP.Addr)
		if err != nil {
			return fmt.Errorf("validate http address: %w", err)
		}

		if !isLocal {
			return fmt.Errorf("http auth token is required for non-local write-capable http transport")
		}
	}

	return nil
}

func (cfg *Config) validateTelemetry() error {
	if !cfg.Telemetry.Enabled {
		return nil
	}

	if cfg.Telemetry.ReadTimeout < 0 {
		return fmt.Errorf("telemetry read timeout must not be negative")
	}

	if cfg.Telemetry.IDLETimeout < 0 {
		return fmt.Errorf("telemetry idle timeout must not be negative")
	}

	_, _, err := net.SplitHostPort(cfg.Telemetry.Addr)
	if err != nil {
		return fmt.Errorf("validate telemetry address: %w", err)
	}

	return nil
}

func LoadFromEnv(ctx context.Context) (*Config, error) {
	var cfg Config

	err := env.ProcessWith(ctx, &env.Config{
		Target:   &cfg,
		Lookuper: env.PrefixLookuper(prefix, env.OsLookuper()),
	})
	if err != nil {
		return nil, fmt.Errorf("process config from env: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

// isLocalHTTPAddr reports whether an HTTP listen address is bound to a
// loopback or localhost interface, making it safe to allow unauthenticated
// write-capable access.
func isLocalHTTPAddr(addr string) (bool, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false, fmt.Errorf("parse http address: %w", err)
	}

	if host == "localhost" {
		return true, nil
	}

	ip := net.ParseIP(host)
	if ip != nil {
		return ip.IsLoopback(), nil
	}

	return false, nil
}
