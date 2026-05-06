package volume

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetters(t *testing.T) {
	t.Run("create volume params", func(t *testing.T) {
		params := NewCreateVolumeParams().
			SetName("  vol-1  ").
			SetDriver("  local  ").
			SetDriverOpts(map[string]string{"size": "10Gi"}).
			SetLabels(map[string]string{"env": "test"})

		assert.Equal(t, "vol-1", params.Name)
		assert.Equal(t, "local", params.Driver)
		assert.Equal(t, map[string]string{"size": "10Gi"}, params.DriverOpts)
		assert.Equal(t, map[string]string{"env": "test"}, params.Labels)
	})

	t.Run("remove volume params", func(t *testing.T) {
		params := NewRemoveVolumeParams().SetName("  vol-1  ").SetForce(true)
		assert.Equal(t, "vol-1", params.Name)
		assert.True(t, params.Force)
	})

	t.Run("list volumes params", func(t *testing.T) {
		params := NewListVolumesParams().SetDangling(true)
		assert.True(t, params.Dangling)
	})

	t.Run("inspect volume params", func(t *testing.T) {
		params := NewInspectVolumeParams().SetName("  vol-1  ")
		assert.Equal(t, "vol-1", params.Name)
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
