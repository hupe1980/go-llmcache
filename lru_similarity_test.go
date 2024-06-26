package llmcache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUSimilarityEngine_LookupAndUpdate(t *testing.T) {
	// Create a test embedder implementation for the LRUSimilarityEngine.
	mockEmbedder := &mockEmbedder{
		embeddings: map[string][]float32{
			"prompt1": {0.1, 0.2, 0.3, 0.4},
			"prompt2": {0.2, 0.2, 0.3, 0.4},
			"prompt3": {-0.1, -0.2, -0.3, -0.4},
		},
	}

	t.Run("Lookup", func(t *testing.T) {
		// Create a new LRUSimilarityEngine instance with the test embedder.
		cache, err := NewLRUSimilarityEngine[string](mockEmbedder)
		assert.NoError(t, err)

		t.Run("Hit", func(t *testing.T) {
			ctx := context.TODO()
			prompt := "prompt1"
			result := "result1"

			err = cache.Update(ctx, prompt, result)
			assert.NoError(t, err)

			foundResult, ok := cache.Lookup(ctx, prompt)
			assert.True(t, ok)
			assert.Equal(t, result, foundResult)
		})

		t.Run("Similarity Hit", func(t *testing.T) {
			ctx := context.TODO()
			prompt := "prompt1"
			result := "result1"

			err = cache.Update(ctx, prompt, result)
			assert.NoError(t, err)

			similarityPrompt := "prompt2"
			foundResult, ok := cache.Lookup(ctx, similarityPrompt)
			assert.True(t, ok)
			assert.Equal(t, result, foundResult)
		})

		t.Run("Miss", func(t *testing.T) {
			otherPrompt := "prompt3"
			foundResult, ok := cache.Lookup(context.TODO(), otherPrompt)
			assert.False(t, ok)
			assert.Equal(t, "", foundResult)
		})
	})
}

// mockEmbedder is a mock implementation of the Embedder interface for testing.
type mockEmbedder struct {
	embeddings map[string][]float32
}

// EmbedQuery is a mock implementation of the Embedder's EmbedQuery method.
func (e *mockEmbedder) EmbedText(ctx context.Context, text string) ([]float32, error) {
	// Mock implementation logic goes here.
	// Return the embedding vector for the given prompt.
	return e.embeddings[text], nil
}
