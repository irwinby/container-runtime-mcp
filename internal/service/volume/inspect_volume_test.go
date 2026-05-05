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

func TestServiceInspectVolume(t *testing.T) {
	type given struct {
		params InspectVolumeParams
		result providers.VolumeInspect
		err    error
	}

	type want struct {
		called bool
		name   string
		result VolumeInspect
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewInspectVolumeParams().SetName("vol1"), result: providers.NewVolumeInspect().SetName("vol1").SetDriver("local")},
			want:  want{called: true, name: "vol1", result: VolumeInspect{Name: "vol1", Driver: "local"}},
		},
		"trimmed name": {
			given: given{params: NewInspectVolumeParams().SetName("  vol1  "), result: providers.NewVolumeInspect().SetName("vol1")},
			want:  want{called: true, name: "vol1", result: VolumeInspect{Name: "vol1"}},
		},
		"empty name": {
			given: given{params: NewInspectVolumeParams().SetName("")},
			want:  want{},
		},
		"whitespace name": {
			given: given{params: NewInspectVolumeParams().SetName("   ")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewInspectVolumeParams().SetName("vol1"), err: errors.New("docker error")},
			want:  want{called: true, name: "vol1"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := volumemock.NewMockProviderClient(t)

			if tt.want.called {
				mockClient.On("InspectVolume", mock.Anything, providers.InspectVolumeParams{
					Name: tt.want.name,
				}).Return(tt.given.result, tt.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.InspectVolume(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			if !tt.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}
