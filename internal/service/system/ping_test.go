package system

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	systemmock "github.com/irwinby/container-runtime-mcp/internal/service/system/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServicePing(t *testing.T) {
	type given struct {
		result providers.PingResult
		err    error
	}

	type want struct {
		result PingResult
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				result: providers.PingResult{APIVersion: "1.43"},
			},
			want: want{
				result: PingResult{APIVersion: "1.43"},
			},
		},
		"error": {
			given: given{
				err: errors.New("provider error"),
			},
			want: want{
				result: PingResult{},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := systemmock.NewMockProviderClient(t)

			mockClient.On("Ping", mock.Anything).Return(tt.given.result, tt.given.err).Once()

			service := NewService(mockClient, zap.NewNop())

			got, err := service.Ping(context.Background())

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.result, got)
		})
	}
}
