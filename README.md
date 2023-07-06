# 🧠 go-llmcache
![Build Status](https://github.com/hupe1980/go-llmcache/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/go-llmcache.svg)](https://pkg.go.dev/github.com/hupe1980/go-llmcache)
> go-llmcache is a Go package that provides a cache implementation for storing and retrieving results of language model (LLM) requests. It utilizes an LRU (Least Recently Used) cache strategy for efficient management of cached entries. The cache is designed to work with LLM requests, where each request is associated with an embedding vector. The package includes functionality to calculate the cosine similarity between two embedding vectors.

## Features
- Caching of LLM request results for fast retrieval
- LRU cache strategy for efficient management of cached entries
- Calculation of cosine similarity between embedding vectors
- Simple and easy-to-use API

## Installation
Use Go modules to include go-llmcache in your project:
```bash
go get github.com/hupe1980/go-llmcache
```

## Usage
```go
import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hupe1980/golc/embedding"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
)

func main() {
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
}
```
Output:
```text
Result: 1879
Result(*** HIT ***): 1879
Result(*** HIT ***): 1879
Result: 1912
```

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

## Contributing
Contributions are welcome! Feel free to open an issue or submit a pull request for any improvements or new features you would like to see.