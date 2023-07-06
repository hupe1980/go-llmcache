package llmcache

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hupe1980/golc/embedding"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
)

func ExampleLRUSimilarityEngine() {
	apiKey := os.Getenv("OPENAI_API_KEY")

	embedder, _ := embedding.NewOpenAI(apiKey)
	openai, _ := llm.NewOpenAI(apiKey)
	engine, _ := NewLRUSimilarityEngine[*schema.ModelResult](embedder)

	cache := New[*schema.ModelResult](engine)

	ctx := context.Background()

	prompts := []string{
		"What year was Einstein born? Return only the year!",
		"What year was Albert Einstein born? Return only the year!",
		"In what year was albert einstein born? Return only the year!",
		"What year was Alan Turing born? Return only the year!",
	}

	for _, prompt := range prompts {
		if result, ok := cache.Lookup(ctx, prompt); ok {
			fmt.Println("Result(*** HIT ***):", strings.ReplaceAll(result.Generations[0].Text, "\n", ""))
			continue
		}
		// If no similar result found in cache, perform the actual LLM lookup
		result, _ := openai.Generate(ctx, prompt)

		_ = cache.Update(ctx, prompt, result)

		fmt.Println("Result:", strings.ReplaceAll(result.Generations[0].Text, "\n", ""))
	}

	// Expected Output:
	// Result: 1879
	// Result(*** HIT ***): 1879
	// Result(*** HIT ***): 1879
	// Result: 1912
}
