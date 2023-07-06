package llmcache

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
)

// Compile time check to ensure LRUSimilarityEngine satisfies the Engine interface.
var _ Engine[any] = (*LRUSimilarityEngine[any])(nil)

// LRUSimilarityEngineOptions contains options for configuring the LRUSimilarityEngine.
type LRUSimilarityEngineOptions struct {
	// Inherits options from LRUEngine.
	LRUEngineOptions
	// Threshold is the minimum cosine similarity required for a result to be considered a match.
	Threshold float64
}

// LRUSimilarityEngine is a cache engine implementation based on LRU (Least Recently Used) strategy
// with cosine similarity matching capability.
type LRUSimilarityEngine[T comparable] struct {
	// threshold is the minimum cosine similarity required for a result to be considered a match.
	threshold float64
	// embedder is the embedding functionality used for similarity calculations.
	embedder Embedder
	// cache is the underlying LRU cache for storing prompt embeddings and results.
	cache *lru.Cache[string, *CacheEntry[T]]
}

// NewLRUSimilarityEngine creates a new LRUSimilarityEngine instance with the provided embedder and options.
// It returns an error if the cache creation fails.
func NewLRUSimilarityEngine[T comparable](embedder Embedder, optFns ...func(o *LRUSimilarityEngineOptions)) (*LRUSimilarityEngine[T], error) {
	opts := LRUSimilarityEngineOptions{
		LRUEngineOptions: LRUEngineOptions{
			MaxCacheSize: 1000, // Default maximum cache size is set to 1000 entries.
		},
		Threshold: 0.95, // Default threshold is set to 0.95.
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	cache, err := lru.New[string, *CacheEntry[T]](opts.MaxCacheSize)
	if err != nil {
		return nil, err
	}

	return &LRUSimilarityEngine[T]{
		threshold: opts.Threshold,
		embedder:  embedder,
		cache:     cache,
	}, nil
}

// Lookup retrieves the most similar cached result associated with the given prompt.
// It returns the result and a boolean indicating whether a match was found.
func (e *LRUSimilarityEngine[T]) Lookup(ctx context.Context, prompt string) (T, bool) {
	if entry, ok := e.cache.Get(prompt); ok {
		return entry.Result, true
	}

	embedding, err := e.embedder.EmbedQuery(ctx, prompt)
	if err != nil {
		return *new(T), false
	}

	var (
		maxResult     T
		maxSimilarity float64
	)

	for _, entry := range e.cache.Values() {
		if entry.Result == *new(T) {
			continue
		}

		otherEmbedding := entry.Embedding
		similarity := Similarity(embedding, otherEmbedding)

		if similarity > maxSimilarity {
			maxSimilarity = similarity
			maxResult = entry.Result
		}
	}

	if maxSimilarity > e.threshold {
		return maxResult, true
	}

	// Store the embedding in the cache
	e.cache.Add(prompt, &CacheEntry[T]{
		Embedding: embedding,
		//Result:    nil,
	})

	return *new(T), false
}

// Update updates the cache with the provided prompt and result.
// It retrieves the embedding if available, or embeds the prompt if it is a new entry.
func (e *LRUSimilarityEngine[T]) Update(ctx context.Context, prompt string, result T) error {
	var embedding []float64

	if entry, ok := e.cache.Get(prompt); ok {
		embedding = entry.Embedding
	} else {
		e, err := e.embedder.EmbedQuery(ctx, prompt)
		if err != nil {
			return err
		}

		embedding = e
	}

	// Store the result in the cache
	e.cache.Add(prompt, &CacheEntry[T]{
		Embedding: embedding,
		Result:    result,
	})

	return nil
}
