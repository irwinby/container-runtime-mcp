package volume

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	volumemock "github.com/irwinby/container-runtime-mcp/internal/service/volume/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServiceListVolumes(t *testing.T) {
	type given struct {
		params ListVolumesParams
		result []providers.Volume
		err    error
	}

	type want struct {
		called bool
		result []Volume
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				params: NewListVolumesParams().SetDangling(true),
				result: []providers.Volume{providers.NewVolume().SetName("vol1").SetDriver("local")},
			},
			want: want{
				called: true,
				result: []Volume{{Name: "vol1", Driver: "local"}},
			},
		},
		"provider error": {
			given: given{
				params: NewListVolumesParams(),
				err:    errors.New("docker error"),
			},
			want: want{called: true},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := volumemock.NewMockProviderClient(t)

			if tt.want.called {
				mockClient.On("ListVolumes", mock.Anything, providers.ListVolumesParams{
					Dangling: tt.given.params.Dangling,
				}).Return(tt.given.result, tt.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.ListVolumes(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}
