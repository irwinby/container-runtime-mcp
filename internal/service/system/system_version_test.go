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

func TestServiceSystemVersion(t *testing.T) {
	type given struct {
		result providers.SystemVersion
		err    error
	}

	type want struct {
		result SystemVersion
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				result: providers.SystemVersion{Version: "24.0.0"},
			},
			want: want{
				result: SystemVersion{Version: "24.0.0"},
			},
		},
		"error": {
			given: given{
				err: errors.New("provider error"),
			},
			want: want{
				result: SystemVersion{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := systemmock.NewMockProviderClient(t)

			mockClient.On("SystemVersion", mock.Anything).Return(test.given.result, test.given.err).Once()

			service := NewService(mockClient, zap.NewNop())

			got, err := service.SystemVersion(context.Background())

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.result, got)
		})
	}
}
