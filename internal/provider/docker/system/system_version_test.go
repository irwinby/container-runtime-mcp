package system

import (
	"context"
	"errors"
	"testing"

	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderSystemVersion(t *testing.T) {
	type given struct {
		err error
	}

	type want struct {
		version    string
		apiVersion string
		os         string
		arch       string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{},
			want: want{
				version:    "25.0.0",
				apiVersion: "1.45",
				os:         "linux",
				arch:       "arm64",
			},
		},
		"error": {
			given: given{err: errors.New("docker error")},
			want:  want{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ServerVersion", mock.Anything, client.ServerVersionOptions{}).Return(client.ServerVersionResult{
				Version:    test.want.version,
				APIVersion: test.want.apiVersion,
				Os:         test.want.os,
				Arch:       test.want.arch,
			}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.SystemVersion(context.Background())

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.version, result.Version)
			require.Equal(t, test.want.apiVersion, result.APIVersion)
			require.Equal(t, test.want.os, result.Os)
			require.Equal(t, test.want.arch, result.Arch)
		})
	}
}
