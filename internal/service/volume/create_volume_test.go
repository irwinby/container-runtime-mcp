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

func TestServiceCreateVolume(t *testing.T) {
	type given struct {
		params CreateVolumeParams
		result providers.VolumeInspect
		err    error
	}

	type want struct {
		called bool
		result VolumeInspect
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewCreateVolumeParams().SetName("vol1").SetDriver("local"), result: providers.NewVolumeInspect().SetName("vol1").SetDriver("local")},
			want:  want{called: true, result: VolumeInspect{Name: "vol1", Driver: "local"}},
		},
		"provider error": {
			given: given{params: NewCreateVolumeParams().SetName("vol1"), err: errors.New("docker error")},
			want:  want{called: true},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := volumemock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("CreateVolume", mock.Anything, providers.CreateVolumeParams{
					Name:   test.given.params.Name,
					Driver: test.given.params.Driver,
				}).Return(test.given.result, test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.CreateVolume(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.result, result)
		})
	}
}

func TestServiceCreateVolume_ReadOnly(t *testing.T) {
	mockClient := volumemock.NewMockProviderClient(t)

	policy := services.NewPolicy(true)
	service := NewService(mockClient, policy, zap.NewNop())

	_, err := service.CreateVolume(context.Background(), NewCreateVolumeParams().SetName("vol1"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
