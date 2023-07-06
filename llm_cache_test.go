package llmcache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLLMCache_LookupAndUpdate(t *testing.T) {
	// Create a mock engine implementation for the LLMCache.
	mockEngine := &mockEngine[string]{
		cache: map[string]string{
			"prompt1": "result1",
		},
	}

	// Create a new LLMCache instance with the test engine.
	cache := New[string](mockEngine)

	// Define the test cases.
	testCases := []struct {
		name           string
		prompt         string
		expectedResult string
		expectedFound  bool
	}{
		{
			name:           "Existing Entry",
			prompt:         "prompt1",
			expectedResult: "result1",
			expectedFound:  true,
		},
		{
			name:           "Non-existing Entry",
			prompt:         "prompt2",
			expectedResult: "",
			expectedFound:  false,
		},
	}

	// Run each test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Perform the lookup operation.
			result, found := cache.Lookup(context.Background(), tc.prompt)

			// Assert the results.
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedFound, found)
		})
	}
}

// mockEngine is a mock implementation of the Engine interface for testing.
type mockEngine[T any] struct {
	cache map[string]T
}

// Lookup is a mock implementation of the Engine's Lookup method.
func (e *mockEngine[T]) Lookup(ctx context.Context, prompt string) (T, bool) {
	result, found := e.cache[prompt]
	return result, found
}

// Update is a mock implementation of the Engine's Update method.
func (e *mockEngine[T]) Update(ctx context.Context, prompt string, result T) error {
	e.cache[prompt] = result
	return nil
}
