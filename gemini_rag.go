package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiService implements both Embedder and Generator interfaces
type GeminiService struct {
	client *genai.Client
}

func NewGeminiService(ctx context.Context, apiKey string) (*GeminiService, error) {
	// Initialize the Gemini client with the provided API key.
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return &GeminiService{client: client}, nil
}

// Embed implements the Embedder interface.
func (g *GeminiService) Embed(ctx context.Context, text string) ([]float32, error) {
	// The specific embedding model used (e.g., text-embedding-004)
	embModel := g.client.EmbeddingModel("text-embedding-004")

	resp, err := embModel.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}
	// Assuming a single embedding for a single text input
	if len(resp.Embedding.Values) == 0 {
		return nil, fmt.Errorf("received empty embedding")
	}
	return resp.Embedding.Values, nil
}

// Generate implements the Generator interface.
func (g *GeminiService) Generate(ctx context.Context, augmentedPrompt string) (string, error) {
	// The specific generative model used (e.g., gemini-2.5-pro)
	genModel := g.client.GenerativeModel("gemini-2.5-pro")
	resp, err := genModel.GenerateContent(
		ctx,
		genai.Text(augmentedPrompt),
	)
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	// Process the streamed/response content
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
	}

	return "Could not generate a response.", nil
}

// RunRAG orchestrates the entire RAG flow.
func RunRAG(ctx context.Context, query string, retriever Retriever, embedder Embedder, generator Generator) (string, error) {
	// 1. Embed the user's query
	queryVector, err := embedder.Embed(ctx, query)
	if err != nil {
		return "", fmt.Errorf("RAG: query embedding error: %w", err)
	}

	// 2. Retrieve context from the Vector DB
	const topK = 3 // Number of chunks to retrieve
	retrievedChunks, err := retriever.Search(ctx, queryVector, topK)
	if err != nil {
		return "", fmt.Errorf("RAG: retrieval error: %w", err)
	}

	// Check if any context was found
	if len(retrievedChunks) == 0 {
		return "I could not find relevant information in the knowledge base.", nil
	}

	// 3. Augment the prompt
	context := strings.Join(retrievedChunks, "\n---\n")
	augmentedPrompt := fmt.Sprintf(`
        You are an expert assistant. Use ONLY the following context to answer the question. 
        If the answer is not present in the context, state that the information is not available.

        --- CONTEXT ---
        %s

        --- QUESTION ---
        %s
    `, context, query)

	// 4. Generate the final response
	finalResponse, err := generator.Generate(ctx, augmentedPrompt)
	if err != nil {
		return "", fmt.Errorf("RAG: final generation error: %w", err)
	}

	return finalResponse, nil
}
