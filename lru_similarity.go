package llmcache

import (
	"context"
	"math"

	lru "github.com/hashicorp/golang-lru/v2"
)

// Compile time check to ensure LRUSimilarityEngine satisfies the Engine interface.
var _ Engine[any] = (*LRUSimilarityEngine[any])(nil)

// DistanceFunc represents a function for calculating the distance between two vectors
type DistanceFunc func(v1, v2 []float32) (float32, error)

// LRUSimilarityEngineOptions contains options for configuring the LRUSimilarityEngine.
type LRUSimilarityEngineOptions struct {
	// Inherits options from LRUEngine.
	LRUEngineOptions
	// DistanceFunc represents the distance function used for calculating the similarity between embeddings.
	DistanceFunc DistanceFunc
	// Threshold is the maximum distance allowed for a result to be considered a match.
	Threshold float32
	// ReturnFirst is a boolean flag indicating whether to return the first match found during lookup.
	// If set to true, the engine will return the first match found within the threshold distance.
	ReturnFirst bool
}

// LRUSimilarityEngine is a cache engine implementation based on LRU (Least Recently Used) strategy
// with cosine similarity matching capability.
type LRUSimilarityEngine[T comparable] struct {
	// embedder is the embedding functionality used for similarity calculations.
	embedder Embedder
	// cache is the underlying LRU cache for storing prompt embeddings and results.
	cache *lru.Cache[string, *CacheEntry[T]]
	// opts contains options for configuring the LRUSimilarityEngine
	opts LRUSimilarityEngineOptions
}

// NewLRUSimilarityEngine creates a new LRUSimilarityEngine instance with the provided embedder and options.
// It returns an error if the cache creation fails.
func NewLRUSimilarityEngine[T comparable](embedder Embedder, optFns ...func(o *LRUSimilarityEngineOptions)) (*LRUSimilarityEngine[T], error) {
	opts := LRUSimilarityEngineOptions{
		LRUEngineOptions: LRUEngineOptions{
			MaxCacheSize: 1000,
		},
		DistanceFunc: SquaredL2,
		Threshold:    float32(0.50),
		ReturnFirst:  false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	cache, err := lru.New[string, *CacheEntry[T]](opts.MaxCacheSize)
	if err != nil {
		return nil, err
	}

	return &LRUSimilarityEngine[T]{
		embedder: embedder,
		cache:    cache,
		opts:     opts,
	}, nil
}

// Lookup retrieves the most similar cached result associated with the given text.
// It returns the result and a boolean indicating whether a match was found.
func (e *LRUSimilarityEngine[T]) Lookup(ctx context.Context, text string) (T, bool) {
	if entry, ok := e.cache.Get(text); ok {
		return entry.Result, true
	}

	embedding, err := e.embedder.EmbedText(ctx, text)
	if err != nil {
		return *new(T), false
	}

	var (
		result T
	)

	found := false
	minDistance := float32(math.MaxFloat32)

	for _, entry := range e.cache.Values() {
		if entry.Result == *new(T) {
			continue
		}

		otherEmbedding := entry.Embedding

		distance, err := e.opts.DistanceFunc(embedding, otherEmbedding)
		if err != nil {
			return *new(T), false
		}

		if distance < e.opts.Threshold && distance < minDistance {
			minDistance = distance
			result = entry.Result
			found = true

			if e.opts.ReturnFirst {
				return result, true
			}
		}
	}

	if found {
		return result, true
	}

	// Store the embedding in the cache
	e.cache.Add(text, &CacheEntry[T]{
		Embedding: embedding,
		Result:    *new(T),
	})

	return *new(T), false
}

// Update updates the cache with the provided prompt and result.
// It retrieves the embedding if available, or embeds the prompt if it is a new entry.
func (e *LRUSimilarityEngine[T]) Update(ctx context.Context, prompt string, result T) error {
	if entry, ok := e.cache.Get(prompt); ok {
		if entry.Result == result {
			return nil // nothing to do
		}

		e.cache.Remove(prompt)

		e.cache.Add(prompt, &CacheEntry[T]{
			Embedding: entry.Embedding,
			Result:    result,
		})
	} else {
		embedding, err := e.embedder.EmbedText(ctx, prompt)
		if err != nil {
			return err
		}

		e.cache.Add(prompt, &CacheEntry[T]{
			Embedding: embedding,
			Result:    result,
		})
	}

	return nil
}

// Clear clears the cache, removing all entries.
// It returns an error if the clear operation fails.
func (e *LRUSimilarityEngine[T]) Clear(ctx context.Context) error {
	e.cache.Purge()
	return nil
}
