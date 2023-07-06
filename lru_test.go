package llmcache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUEngine(t *testing.T) {
	t.Run("Lookup", func(t *testing.T) {
		engine, err := NewLRUEngine[int]()
		assert.NoError(t, err)
		t.Run("Hit", func(t *testing.T) {
			ctx := context.TODO()
			prompt := "Hello, World!"
			result := 42

			err = engine.Update(ctx, prompt, result)
			assert.NoError(t, err)
			foundResult, ok := engine.Lookup(ctx, prompt)
			assert.True(t, ok)
			assert.Equal(t, result, foundResult)
		})

		t.Run("Miss", func(t *testing.T) {
			otherPrompt := "Goodbye"
			foundResult, ok := engine.Lookup(context.TODO(), otherPrompt)
			assert.False(t, ok)
			assert.Equal(t, foundResult, 0)
		})
	})

	t.Run("Update", func(t *testing.T) {
		engine, err := NewLRUEngine[int]()
		assert.NoError(t, err)

		ctx := context.TODO()
		prompt := "Hello, World!"
		result := 42

		err = engine.Update(ctx, prompt, result)
		assert.NoError(t, err)

		// Verify that the result is stored in the cache
		foundResult, ok := engine.Lookup(ctx, prompt)
		assert.True(t, ok)
		assert.Equal(t, result, foundResult)
	})

	t.Run("Clear", func(t *testing.T) {
		engine, err := NewLRUEngine[int]()
		assert.NoError(t, err)

		ctx := context.TODO()
		prompt := "Hello, World!"
		result := 42

		err = engine.Update(ctx, prompt, result)
		assert.NoError(t, err)

		_ = engine.Clear(ctx)

		// Verify that the result is not stored in the cache
		foundResult, ok := engine.Lookup(ctx, prompt)
		assert.False(t, ok)
		assert.Equal(t, 0, foundResult)
	})
}
