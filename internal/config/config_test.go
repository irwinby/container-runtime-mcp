package config

import (
	"context"
	"testing"
	"time"

	"github.com/irwinby/container-runtime-mcp/internal/testing/env"
	"github.com/irwinby/container-runtime-mcp/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	type given struct {
		env map[string]string
	}

	type want struct {
		cfg *Config
		err bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"defaults": {
			given: given{
				env: map[string]string{},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportStdio,
							HTTP: &HTTPTransportConfig{
								Addr:           "127.0.0.1:8080",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"custom values": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_SERVER_NAME":              "CustomRuntime",
					"CONTAINER_RUNTIME_MCP_SERVER_TITLE":             "Custom Title",
					"CONTAINER_RUNTIME_MCP_SERVER_VERSION":           "2.0.0",
					"CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT": "5m",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "CustomRuntime",
						Title:   "Custom Title",
						Version: "2.0.0",
						TransportConfig: TransportConfig{
							Type: TransportStdio,
							HTTP: &HTTPTransportConfig{
								Addr:           "127.0.0.1:8080",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 5 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"negative timeout": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT": "-5m",
				},
			},
			want: want{
				err: true,
			},
		},
		"http transport local": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT":            "http",
					"CONTAINER_RUNTIME_MCP_HTTP_ADDR":            "127.0.0.1:3000",
					"CONTAINER_RUNTIME_MCP_HTTP_PATH":            "/runtime",
					"CONTAINER_RUNTIME_MCP_HTTP_SESSION_TIMEOUT": "1h",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportHTTP,
							HTTP: &HTTPTransportConfig{
								Addr:           "127.0.0.1:3000",
								Path:           "/runtime",
								SessionTimeout: 1 * time.Hour,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"http transport with auth": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT":            "http",
					"CONTAINER_RUNTIME_MCP_HTTP_ADDR":            "0.0.0.0:3000",
					"CONTAINER_RUNTIME_MCP_HTTP_PATH":            "/runtime",
					"CONTAINER_RUNTIME_MCP_HTTP_SESSION_TIMEOUT": "1h",
					"CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN":      "my-secret-token",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportHTTP,
							HTTP: &HTTPTransportConfig{
								Addr:           "0.0.0.0:3000",
								Path:           "/runtime",
								SessionTimeout: 1 * time.Hour,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "my-secret-token",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"read only mode": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_READ_ONLY": "true",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportStdio,
							HTTP: &HTTPTransportConfig{
								Addr:           "127.0.0.1:8080",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					ReadOnly:               true,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"invalid transport": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT": "ws",
				},
			},
			want: want{
				err: true,
			},
		},
		"http path missing leading slash": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT": "http",
					"CONTAINER_RUNTIME_MCP_HTTP_PATH": "mcp",
				},
			},
			want: want{
				err: true,
			},
		},
		"negative http session timeout": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT":            "http",
					"CONTAINER_RUNTIME_MCP_HTTP_SESSION_TIMEOUT": "-5m",
				},
			},
			want: want{
				err: true,
			},
		},
		"negative http read timeout": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT":         "http",
					"CONTAINER_RUNTIME_MCP_HTTP_READ_TIMEOUT": "-5s",
				},
			},
			want: want{
				err: true,
			},
		},
		"negative http idle timeout": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT":         "http",
					"CONTAINER_RUNTIME_MCP_HTTP_IDLE_TIMEOUT": "-5s",
				},
			},
			want: want{
				err: true,
			},
		},
		"non-local http without auth": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT": "http",
					"CONTAINER_RUNTIME_MCP_HTTP_ADDR": "0.0.0.0:3000",
				},
			},
			want: want{
				err: true,
			},
		},
		"non-local http with auth": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT":       "http",
					"CONTAINER_RUNTIME_MCP_HTTP_ADDR":       "0.0.0.0:3000",
					"CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN": "secret",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportHTTP,
							HTTP: &HTTPTransportConfig{
								Addr:           "0.0.0.0:3000",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "secret",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"non-local http read-only without auth": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT": "http",
					"CONTAINER_RUNTIME_MCP_HTTP_ADDR": "0.0.0.0:3000",
					"CONTAINER_RUNTIME_MCP_READ_ONLY": "true",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportHTTP,
							HTTP: &HTTPTransportConfig{
								Addr:           "0.0.0.0:3000",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					ReadOnly:               true,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"localhost http without auth": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT": "http",
					"CONTAINER_RUNTIME_MCP_HTTP_ADDR": "localhost:3000",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportHTTP,
							HTTP: &HTTPTransportConfig{
								Addr:           "localhost:3000",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"invalid http address": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT": "http",
					"CONTAINER_RUNTIME_MCP_HTTP_ADDR": "not-an-address",
				},
			},
			want: want{
				err: true,
			},
		},
		"custom log level": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_LOG_LEVEL": "debug",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportStdio,
							HTTP: &HTTPTransportConfig{
								Addr:           "127.0.0.1:8080",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.DebugLevel,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"invalid log level": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_LOG_LEVEL": "trace",
				},
			},
			want: want{
				err: true,
			},
		},
		"negative remote operation timeout validate": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT": "-1m",
				},
			},
			want: want{
				err: true,
			},
		},
		"invalid duration parse": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT": "not-a-duration",
				},
			},
			want: want{
				err: true,
			},
		},
		"non-loopback non-localhost address": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TRANSPORT":       "http",
					"CONTAINER_RUNTIME_MCP_HTTP_ADDR":       "192.168.1.1:3000",
					"CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN": "secret",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportHTTP,
							HTTP: &HTTPTransportConfig{
								Addr:           "192.168.1.1:3000",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "secret",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      false,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"telemetry enabled with defaults": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TELEMETRY_ENABLED": "true",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportStdio,
							HTTP: &HTTPTransportConfig{
								Addr:           "127.0.0.1:8080",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      true,
						Addr:         "127.0.0.1:9090",
						PPROFEnabled: false,
						ReadTimeout:  10 * time.Second,
						IDLETimeout:  120 * time.Second,
					},
				},
			},
		},
		"telemetry enabled custom values": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TELEMETRY_ENABLED":       "true",
					"CONTAINER_RUNTIME_MCP_TELEMETRY_ADDR":          "127.0.0.1:9091",
					"CONTAINER_RUNTIME_MCP_TELEMETRY_PPROF_ENABLED": "true",
					"CONTAINER_RUNTIME_MCP_TELEMETRY_READ_TIMEOUT":  "5s",
					"CONTAINER_RUNTIME_MCP_TELEMETRY_IDLE_TIMEOUT":  "60s",
				},
			},
			want: want{
				cfg: &Config{
					MCPServer: MCPServer{
						Name:    "Container Runtime",
						Title:   "",
						Version: "1.0.0",
						TransportConfig: TransportConfig{
							Type: TransportStdio,
							HTTP: &HTTPTransportConfig{
								Addr:           "127.0.0.1:8080",
								Path:           "/mcp",
								SessionTimeout: 30 * time.Minute,
								ReadTimeout:    10 * time.Second,
								IDLETimeout:    120 * time.Second,
								AuthToken:      "",
							},
						},
					},
					RemoteOperationTimeout: 10 * time.Minute,
					LogLevel:               logger.InfoLevel,
					Telemetry: TelemetryConfig{
						Enabled:      true,
						Addr:         "127.0.0.1:9091",
						PPROFEnabled: true,
						ReadTimeout:  5 * time.Second,
						IDLETimeout:  60 * time.Second,
					},
				},
			},
		},
		"telemetry negative read timeout": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TELEMETRY_ENABLED":      "true",
					"CONTAINER_RUNTIME_MCP_TELEMETRY_READ_TIMEOUT": "-5s",
				},
			},
			want: want{
				err: true,
			},
		},
		"telemetry negative idle timeout": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TELEMETRY_ENABLED":      "true",
					"CONTAINER_RUNTIME_MCP_TELEMETRY_IDLE_TIMEOUT": "-5s",
				},
			},
			want: want{
				err: true,
			},
		},
		"telemetry invalid address": {
			given: given{
				env: map[string]string{
					"CONTAINER_RUNTIME_MCP_TELEMETRY_ENABLED": "true",
					"CONTAINER_RUNTIME_MCP_TELEMETRY_ADDR":    "not-an-address",
				},
			},
			want: want{
				err: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Clear all relevant env vars before each test case
			// to avoid host environment leaking into assertions.
			for _, k := range []string{
				"CONTAINER_RUNTIME_MCP_SERVER_NAME",
				"CONTAINER_RUNTIME_MCP_SERVER_TITLE",
				"CONTAINER_RUNTIME_MCP_SERVER_VERSION",
				"CONTAINER_RUNTIME_MCP_TRANSPORT",
				"CONTAINER_RUNTIME_MCP_HTTP_ADDR",
				"CONTAINER_RUNTIME_MCP_HTTP_PATH",
				"CONTAINER_RUNTIME_MCP_HTTP_SESSION_TIMEOUT",
				"CONTAINER_RUNTIME_MCP_HTTP_READ_TIMEOUT",
				"CONTAINER_RUNTIME_MCP_HTTP_IDLE_TIMEOUT",
				"CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN",
				"CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT",
				"CONTAINER_RUNTIME_MCP_READ_ONLY",
				"CONTAINER_RUNTIME_LOG_LEVEL",
				"CONTAINER_RUNTIME_MCP_TELEMETRY_ENABLED",
				"CONTAINER_RUNTIME_MCP_TELEMETRY_ADDR",
				"CONTAINER_RUNTIME_MCP_TELEMETRY_PPROF_ENABLED",
				"CONTAINER_RUNTIME_MCP_TELEMETRY_READ_TIMEOUT",
				"CONTAINER_RUNTIME_MCP_TELEMETRY_IDLE_TIMEOUT",
			} {
				env.Unset(t, k)
			}

			for k, v := range test.given.env {
				t.Setenv(k, v)
			}

			cfg, err := LoadFromEnv(context.Background())

			if test.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.cfg, cfg)
		})
	}
}

func TestLoadFromEnv_Prefix(t *testing.T) {
	// Ensure prefix constant includes the separator
	require.Equal(t, "CONTAINER_RUNTIME_", prefix)
}

func TestValidate_NegativeRemoteOperationTimeout(t *testing.T) {
	cfg := &Config{
		MCPServer: MCPServer{
			TransportConfig: TransportConfig{Type: TransportStdio},
		},
		RemoteOperationTimeout: -1 * time.Minute,
		LogLevel:               "info",
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "remote operation timeout must not be negative")
}

func TestIsLocalHTTPAddr(t *testing.T) {
	type given struct {
		addr string
	}

	type want struct {
		local bool
		err   bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"127.0.0.1": {
			given: given{addr: "127.0.0.1:8080"},
			want:  want{local: true},
		},
		"localhost": {
			given: given{addr: "localhost:8080"},
			want:  want{local: true},
		},
		"0.0.0.0": {
			given: given{addr: "0.0.0.0:8080"},
			want:  want{local: false},
		},
		"192.168.1.1": {
			given: given{addr: "192.168.1.1:8080"},
			want:  want{local: false},
		},
		"invalid": {
			given: given{addr: "not-an-address"},
			want:  want{err: true},
		},
		"ipv6 loopback": {
			given: given{addr: "[::1]:8080"},
			want:  want{local: true},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			local, err := isLocalHTTPAddr(test.given.addr)

			if test.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.local, local)
		})
	}
}

func TestLoadFromEnv_UnsetEnv(t *testing.T) {
	for _, key := range []string{
		"CONTAINER_RUNTIME_MCP_SERVER_NAME",
		"CONTAINER_RUNTIME_MCP_SERVER_TITLE",
		"CONTAINER_RUNTIME_MCP_SERVER_VERSION",
		"CONTAINER_RUNTIME_MCP_TRANSPORT",
		"CONTAINER_RUNTIME_MCP_HTTP_ADDR",
		"CONTAINER_RUNTIME_MCP_HTTP_PATH",
		"CONTAINER_RUNTIME_MCP_HTTP_SESSION_TIMEOUT",
		"CONTAINER_RUNTIME_MCP_HTTP_READ_TIMEOUT",
		"CONTAINER_RUNTIME_MCP_HTTP_IDLE_TIMEOUT",
		"CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN",
		"CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT",
		"CONTAINER_RUNTIME_MCP_READ_ONLY",
		"CONTAINER_RUNTIME_LOG_LEVEL",
		"CONTAINER_RUNTIME_MCP_TELEMETRY_ENABLED",
		"CONTAINER_RUNTIME_MCP_TELEMETRY_ADDR",
		"CONTAINER_RUNTIME_MCP_TELEMETRY_PPROF_ENABLED",
		"CONTAINER_RUNTIME_MCP_TELEMETRY_READ_TIMEOUT",
		"CONTAINER_RUNTIME_MCP_TELEMETRY_IDLE_TIMEOUT",
	} {
		env.Unset(t, key)
	}

	cfg, err := LoadFromEnv(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "Container Runtime", cfg.Name)
	assert.Equal(t, "", cfg.Title)
	assert.Equal(t, "1.0.0", cfg.Version)
	assert.Equal(t, TransportStdio, cfg.TransportConfig.Type)
	assert.Equal(t, 10*time.Minute, cfg.RemoteOperationTimeout)
	assert.Equal(t, false, cfg.Telemetry.Enabled)
	assert.Equal(t, "127.0.0.1:9090", cfg.Telemetry.Addr)
	assert.Equal(t, false, cfg.Telemetry.PPROFEnabled)
	assert.Equal(t, 10*time.Second, cfg.Telemetry.ReadTimeout)
	assert.Equal(t, 120*time.Second, cfg.Telemetry.IDLETimeout)
}
