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

func TestHandlerSystemInfo(t *testing.T) {
	type given struct {
		result system.SystemInfo
		err    error
	}

	type want struct {
		called bool
		info   SystemInfo
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				result: system.SystemInfo{ID: "abc", Containers: 5},
			},
			want: want{
				called: true,
				info:   SystemInfo{ID: "abc", Containers: 5},
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
			mockSvc.On("SystemInfo", mock.Anything).Return(tt.given.result, tt.given.err).Once()

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.SystemInfo(context.Background(), nil, SystemInfoInput{})

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.info, output.Info)
		})
	}
}
