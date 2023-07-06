package llmcache

import (
	"context"
)

// Engine is an interface for performing lookup and update operations.
type Engine[T comparable] interface {
	// Lookup retrieves the cached result associated with the given prompt.
	// It returns the result and a boolean indicating whether the result was found.
	Lookup(ctx context.Context, prompt string) (T, bool)

	// Update updates the cache with the provided prompt and result.
	// It returns an error if the update operation fails.
	Update(ctx context.Context, prompt string, result T) error
}

// CacheEntry represents an entry in the cache.
type CacheEntry[T comparable] struct {
	// Embedding is the vector representation of the prompt.
	Embedding []float64
	// Result is the cached result associated with the prompt.
	Result T
}

// Embedder is an interface for embedding queries.
type Embedder interface {
	// EmbedQuery embeds the given text and returns the embedding vector.
	// It returns an error if the embedding operation fails.
	EmbedQuery(ctx context.Context, text string) ([]float64, error)
}

// LLMCache is a cache implementation that utilizes an Engine.
type LLMCache[T comparable] struct {
	// engine is the underlying engine used for lookup and update operations.
	engine Engine[T]
}

// New creates a new LLMCache instance with the provided engine.
func New[T comparable](engine Engine[T]) *LLMCache[T] {
	return &LLMCache[T]{
		engine: engine,
	}
}

// Lookup retrieves the cached result associated with the given prompt.
// It returns the result and a boolean indicating whether the result was found.
func (c *LLMCache[T]) Lookup(ctx context.Context, prompt string) (T, bool) {
	return c.engine.Lookup(ctx, prompt)
}

// Update updates the cache with the provided prompt and result.
// It returns an error if the update operation fails.
func (c *LLMCache[T]) Update(ctx context.Context, prompt string, result T) error {
	return c.engine.Update(ctx, prompt, result)
}
