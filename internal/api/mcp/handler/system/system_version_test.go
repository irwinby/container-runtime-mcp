package system

import (
	"context"
	"errors"
	"testing"

	"github.com/irwinby/container-runtime-mcp/internal/service/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	systemmock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/system/mock"
)

func TestHandlerSystemVersion(t *testing.T) {
	type given struct {
		result system.SystemVersion
		err    error
	}

	type want struct {
		called  bool
		version Version
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				result: system.SystemVersion{Version: "24.0.0", APIVersion: "1.43"},
			},
			want: want{
				called:  true,
				version: Version{Version: "24.0.0", APIVersion: "1.43"},
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

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := systemmock.NewMockSystemService(t)

			mockService.On("SystemVersion", mock.Anything).Return(test.given.result, test.given.err).Once()

			handler := NewToolsHandler(mockService)

			_, output, err := handler.Version(context.Background(), nil, VersionInput{})

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.version, output.Version)
		})
	}
}
