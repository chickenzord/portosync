package version

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	info := Get()

	assert.NotEmpty(t, info.Version)
	assert.NotEmpty(t, info.GitCommit)
	assert.NotEmpty(t, info.BuildDate)
	assert.Equal(t, runtime.Version(), info.GoVersion)
	assert.Contains(t, info.Platform, runtime.GOOS)
	assert.Contains(t, info.Platform, runtime.GOARCH)
}

func TestInfoString(t *testing.T) {
	info := Info{
		Version:   "v1.0.0",
		GitCommit: "abc123",
		BuildDate: "2024-01-01",
		GoVersion: "go1.21.0",
		Platform:  "linux/amd64",
	}

	result := info.String()

	assert.Contains(t, result, "portosync")
	assert.Contains(t, result, "v1.0.0")
	assert.Contains(t, result, "abc123")
	assert.Contains(t, result, "2024-01-01")
	assert.Contains(t, result, "go1.21.0")
	assert.Contains(t, result, "linux/amd64")
}

func TestDefaultValues(t *testing.T) {
	assert.Equal(t, "dev", Version)
	assert.Equal(t, "unknown", GitCommit)
	assert.Equal(t, "unknown", BuildDate)
}
