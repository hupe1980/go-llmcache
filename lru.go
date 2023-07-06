package llmcache

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
)

// Compile time check to ensure LRUEngine satisfies the Engine interface.
var _ Engine[any] = (*LRUEngine[any])(nil)

// LRUEngineOptions contains options for configuring the LRUEngine.
type LRUEngineOptions struct {
	// MaxCacheSize is the maximum number of entries to be stored in the cache.
	MaxCacheSize int
}

// LRUEngine is a cache engine implementation based on LRU (Least Recently Used) strategy.
type LRUEngine[T comparable] struct {
	// cache is the underlying LRU cache.
	cache *lru.Cache[string, T]
}

// NewLRUEngine creates a new LRUEngine instance with the provided options.
// It returns an error if the cache creation fails.
func NewLRUEngine[T comparable](optFns ...func(o *LRUEngineOptions)) (*LRUEngine[T], error) {
	opts := LRUEngineOptions{
		MaxCacheSize: 1000,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	cache, err := lru.New[string, T](opts.MaxCacheSize)
	if err != nil {
		return nil, err
	}

	return &LRUEngine[T]{
		cache: cache,
	}, nil
}

// Lookup retrieves the cached result associated with the given prompt.
// It returns the result and a boolean indicating whether the result was found.
func (e *LRUEngine[T]) Lookup(ctx context.Context, prompt string) (T, bool) {
	return e.cache.Get(prompt)
}

// Update updates the cache with the provided prompt and result.
// It returns an error if the update operation fails.
func (e *LRUEngine[T]) Update(ctx context.Context, prompt string, result T) error {
	e.cache.Add(prompt, result)
	return nil
}

// Clear clears the cache, removing all entries.
// It returns an error if the clear operation fails.
func (e *LRUEngine[T]) Clear(ctx context.Context) error {
	e.cache.Purge()
	return nil
}
