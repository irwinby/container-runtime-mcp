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

func TestServiceRemoveVolume(t *testing.T) {
	type given struct {
		params RemoveVolumeParams
		err    error
	}

	type want struct {
		called bool
		name   string
		force  bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewRemoveVolumeParams().SetName("vol1").SetForce(true)},
			want:  want{called: true, name: "vol1", force: true},
		},
		"empty name": {
			given: given{params: NewRemoveVolumeParams().SetName("")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewRemoveVolumeParams().SetName("vol1"), err: errors.New("docker error")},
			want:  want{called: true, name: "vol1"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := volumemock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("RemoveVolume", mock.Anything, providers.RemoveVolumeParams{
					Name:  test.want.name,
					Force: test.want.force,
				}).Return(test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			err := service.RemoveVolume(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			if !test.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestServiceRemoveVolume_ReadOnly(t *testing.T) {
	mockClient := volumemock.NewMockProviderClient(t)

	policy := services.NewPolicy(true)
	service := NewService(mockClient, policy, zap.NewNop())

	err := service.RemoveVolume(context.Background(), NewRemoveVolumeParams().SetName("vol1"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
