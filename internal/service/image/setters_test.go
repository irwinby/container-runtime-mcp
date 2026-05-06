package image

import (
	"testing"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
)

func TestSetters(t *testing.T) {
	t.Run("list images params", func(t *testing.T) {
		params := NewListImagesParams().SetAll(true).SetSharedSize(true)
		assert.True(t, params.All)
		assert.True(t, params.SharedSize)
	})

	t.Run("inspect image params", func(t *testing.T) {
		params := NewInspectImageParams().SetRef("  nginx:latest  ")
		assert.Equal(t, "nginx:latest", params.Ref)
	})

	t.Run("pull image params", func(t *testing.T) {
		platform := &ocispec.Platform{OS: "linux", Architecture: "amd64"}
		params := NewPullImageParams().
			SetRef("  nginx:latest  ").
			SetAll(true).
			SetPlatform(platform)

		assert.Equal(t, "nginx:latest", params.Ref)
		assert.True(t, params.All)
		assert.Equal(t, platform, params.Platform)
	})

	t.Run("push image params", func(t *testing.T) {
		platform := &ocispec.Platform{OS: "linux", Architecture: "amd64"}
		params := NewPushImageParams().
			SetRef("  nginx:latest  ").
			SetAll(true).
			SetPlatform(platform)

		assert.Equal(t, "nginx:latest", params.Ref)
		assert.True(t, params.All)
		assert.Equal(t, platform, params.Platform)
	})

	t.Run("remove image params", func(t *testing.T) {
		platform := &ocispec.Platform{OS: "linux", Architecture: "amd64"}
		params := NewRemoveImageParams().
			SetRef("  nginx:latest  ").
			SetForce(true).
			SetPruneChildren(true).
			SetPlatform(platform)

		assert.Equal(t, "nginx:latest", params.Ref)
		assert.True(t, params.Force)
		assert.True(t, params.PruneChildren)
		assert.Equal(t, platform, params.Platform)
	})

	t.Run("tag image params", func(t *testing.T) {
		params := NewTagImageParams().SetSource("  nginx  ").SetTarget("  my-nginx  ")
		assert.Equal(t, "nginx", params.Source)
		assert.Equal(t, "my-nginx", params.Target)
	})
}

func TestServiceCanWrite(t *testing.T) {
	t.Run("read only", func(t *testing.T) {
		service := &Service{policy: struct{ ReadOnly bool }{ReadOnly: true}}
		assert.False(t, service.CanWrite())
	})

	t.Run("writable", func(t *testing.T) {
		service := &Service{policy: struct{ ReadOnly bool }{ReadOnly: false}}
		assert.True(t, service.CanWrite())
	})
}
