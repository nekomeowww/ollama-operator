package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOllamaPull(t *testing.T) {
	command := ollamaPull("gemma3:270m")
	assert.Equal(t, "curl -X 'POST' -d '{\"model\":\"gemma3:270m\"}' -H 'Content-Type: application/json' 'http://ollama-models-store:11434/api/pull'", command)
}

func TestOllamaGenerate(t *testing.T) {
	command := ollamaGenerate("gemma3:270m")
	assert.Equal(t, "curl -X 'POST' -d '{\"model\":\"gemma3:270m\"}' -H 'Content-Type: application/json' 'http://ollama-models-store:11434/api/generate'", command)
}
