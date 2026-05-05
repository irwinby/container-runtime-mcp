package system

import (
	"context"
	"testing"

	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderPing(t *testing.T) {
	type given struct {
		err error
	}

	type want struct {
		apiVersion   string
		osType       string
		experimental bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{},
			want: want{
				apiVersion:   "1.45",
				osType:       "linux",
				experimental: true,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("Ping", mock.Anything, client.PingOptions{}).Return(client.PingResult{
				APIVersion:   tt.want.apiVersion,
				OSType:       tt.want.osType,
				Experimental: tt.want.experimental,
			}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)
			result, err := provider.Ping(context.Background())

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.apiVersion, result.APIVersion)
			require.Equal(t, tt.want.osType, result.OSType)
			require.Equal(t, tt.want.experimental, result.Experimental)
		})
	}
}
