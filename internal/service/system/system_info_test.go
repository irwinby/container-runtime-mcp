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

func TestServiceSystemInfo(t *testing.T) {
	type given struct {
		result providers.SystemInfo
		err    error
	}

	type want struct {
		result SystemInfo
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				result: providers.SystemInfo{ID: "abc", Containers: 3},
			},
			want: want{
				result: SystemInfo{ID: "abc", Containers: 3},
			},
		},
		"error": {
			given: given{
				err: errors.New("provider error"),
			},
			want: want{
				result: SystemInfo{},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := systemmock.NewMockProviderClient(t)

			mockClient.On("SystemInfo", mock.Anything).Return(tt.given.result, tt.given.err).Once()

			service := NewService(mockClient, zap.NewNop())

			got, err := service.SystemInfo(context.Background())

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.result, got)
		})
	}
}
