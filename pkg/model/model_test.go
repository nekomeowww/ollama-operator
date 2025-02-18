package model

import (
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelNameFromImage(t *testing.T) {
	ref, err := name.ParseReference(
		"deepseek-r1",
		name.Insecure,
		name.WithDefaultRegistry("https://registry.ollama.ai"),
		name.WithDefaultTag("latest"),
	)
	require.NoError(t, err)
	assert.Equal(t, "deepseek-r1", ref.Context().RepositoryStr())
	assert.Equal(t, "latest", ref.Identifier())
	assert.Equal(t, "deepseek-r1", ref.String())

	ref, err = name.ParseReference(
		"deepseek-r1:7b",
		name.Insecure,
		name.WithDefaultRegistry("https://registry.ollama.ai"),
		name.WithDefaultTag("latest"),
	)
	require.NoError(t, err)
	assert.Equal(t, "deepseek-r1", ref.Context().RepositoryStr())
	assert.Equal(t, "7b", ref.Identifier())
	assert.Equal(t, "deepseek-r1:7b", ref.String())

	ref, err = name.ParseReference(
		"registry.ollama.ai/library/deepseek-r1:7b",
		name.Insecure,
		name.WithDefaultRegistry("https://registry.ollama.ai"),
		name.WithDefaultTag("latest"),
	)
	require.NoError(t, err)
	assert.Equal(t, "library/deepseek-r1", ref.Context().RepositoryStr())
	assert.Equal(t, "deepseek-r1", filepath.Base(ref.Context().RepositoryStr()))
	assert.Equal(t, "7b", ref.Identifier())
	assert.Equal(t, "registry.ollama.ai/library/deepseek-r1:7b", ref.String())
}
