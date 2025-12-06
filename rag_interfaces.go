package main

import (
	"context"
)

// Retriever defines the contract for finding relevant documents (the VDB layer).
type Retriever interface {
	// Search retrieves the top K relevant text chunks based on the query vector.
	Search(ctx context.Context, queryVector []float32, k int) ([]string, error)
}

// Embedder defines the contract for converting text to a vector.
type Embedder interface {
	// Embed converts a single string of text into a vector/embedding.
	Embed(ctx context.Context, text string) ([]float32, error)
}

// Generator defines the contract for generating the final augmented response.
type Generator interface {
	// Generate takes the augmented prompt and returns the final LLM response.
	Generate(ctx context.Context, augmentedPrompt string) (string, error)
}
