package docker

import (
	"context"
	"errors"
	"testing"
	"time"

	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderClose(t *testing.T) {
	type given struct {
		err error
	}

	type want struct {
		err bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"error": {
			given: given{
				err: errors.New("close failed"),
			},
			want: want{
				err: true,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("Close").Return(tt.given.err)

			provider := newProvider(mockClient, 0)
			err := provider.Close()

			if tt.want.err {
				require.Error(t, err)
				require.Equal(t, tt.given.err.Error(), err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestProviderWithTimeout(t *testing.T) {
	type given struct {
		timeout time.Duration
	}

	tests := map[string]struct {
		given        given
		wantDeadline bool
	}{
		"applies timeout": {
			given: given{
				timeout: 100 * time.Millisecond,
			},
			wantDeadline: true,
		},
		"zero timeout disables deadline": {
			given: given{
				timeout: 0,
			},
			wantDeadline: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			provider := newProvider(mockClient, tt.given.timeout)

			ctx, cancel := provider.ContainerProvider.WithTimeout(context.Background())
			defer cancel()

			_, hasDeadline := ctx.Deadline()
			require.Equal(t, tt.wantDeadline, hasDeadline)
		})
	}
}
