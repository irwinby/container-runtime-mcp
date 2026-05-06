package system

import (
	"context"
	"errors"
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
		"error": {
			given: given{err: errors.New("docker error")},
			want:  want{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("Ping", mock.Anything, client.PingOptions{}).Return(client.PingResult{
				APIVersion:   test.want.apiVersion,
				OSType:       test.want.osType,
				Experimental: test.want.experimental,
			}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)
			result, err := provider.Ping(context.Background())

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.apiVersion, result.APIVersion)
			require.Equal(t, test.want.osType, result.OSType)
			require.Equal(t, test.want.experimental, result.Experimental)
		})
	}
}
