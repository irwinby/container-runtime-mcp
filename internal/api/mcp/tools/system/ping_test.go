package system

import (
	"context"
	"errors"
	"testing"

	"github.com/irwinby/container-runtime-mcp/internal/service/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	systemmock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/system/mock"
)

func TestHandlerPing(t *testing.T) {
	type given struct {
		result system.PingResult
		err    error
	}

	type want struct {
		called bool
		ping   PingResult
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				result: system.PingResult{APIVersion: "1.43", OSType: "linux"},
			},
			want: want{
				called: true,
				ping:   PingResult{APIVersion: "1.43", OSType: "linux"},
			},
		},
		"error": {
			given: given{
				err: errors.New("docker error"),
			},
			want: want{
				called: true,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := systemmock.NewMockSystemService(t)
			mockSvc.On("Ping", mock.Anything).Return(tt.given.result, tt.given.err).Once()

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.Ping(context.Background(), nil, PingInput{})

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.ping, output.Ping)
		})
	}
}
