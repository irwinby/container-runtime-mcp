package config

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/irwinby/container-runtime-mcp/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// unsetEnv unsets an environment variable and registers a cleanup
// that restores the original value after the test.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	old, ok := os.LookupEnv(key)
	os.Unsetenv(key)
	t.Cleanup(func() {
		if ok {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	})
}

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
				},
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
			} {
				unsetEnv(t, k)
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
	tests := map[string]struct {
		addr    string
		local   bool
		wantErr bool
	}{
		"127.0.0.1":     {addr: "127.0.0.1:8080", local: true},
		"localhost":     {addr: "localhost:8080", local: true},
		"0.0.0.0":       {addr: "0.0.0.0:8080", local: false},
		"192.168.1.1":   {addr: "192.168.1.1:8080", local: false},
		"invalid":       {addr: "not-an-address", wantErr: true},
		"ipv6 loopback": {addr: "[::1]:8080", local: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			local, err := isLocalHTTPAddr(test.addr)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.local, local)
		})
	}
}

func TestLoadFromEnv_UnsetEnv(t *testing.T) {
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
	} {
		unsetEnv(t, k)
	}

	cfg, err := LoadFromEnv(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Container Runtime", cfg.Name)
	assert.Equal(t, "", cfg.Title)
	assert.Equal(t, "1.0.0", cfg.Version)
	assert.Equal(t, TransportStdio, cfg.TransportConfig.Type)
	assert.Equal(t, 10*time.Minute, cfg.RemoteOperationTimeout)
}
