package llmcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimilarity(t *testing.T) {
	// Define test cases.
	testCases := []struct {
		name       string
		embeddingA []float64
		embeddingB []float64
		expected   float64
	}{
		{
			name:       "Equal Embeddings",
			embeddingA: []float64{0.5, 0.5, 0.5},
			embeddingB: []float64{0.5, 0.5, 0.5},
			expected:   1.0,
		},
		{
			name:       "Orthogonal Embeddings",
			embeddingA: []float64{1, 0, 0},
			embeddingB: []float64{0, 1, 0},
			expected:   0.0,
		},
		{
			name:       "Opposite Embeddings",
			embeddingA: []float64{0.1, 0.2, 0.3},
			embeddingB: []float64{-0.1, -0.2, -0.3},
			expected:   -1.0,
		},
		{
			name:       "Empty Embeddings",
			embeddingA: []float64{},
			embeddingB: []float64{},
			expected:   0.0,
		},
		{
			name:       "Zero Magnitude",
			embeddingA: []float64{0, 0, 0},
			embeddingB: []float64{1, 2, 3},
			expected:   0.0,
		},
	}

	// Run each test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate the similarity.
			result := Similarity(tc.embeddingA, tc.embeddingB)

			// Assert the result.
			assert.InDelta(t, tc.expected, result, 0.0001)
		})
	}
}
