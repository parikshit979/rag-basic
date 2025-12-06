package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// 1. Load environment variables (API Key)
	// Create a .env file in the root directory with GEMINI_API_KEY="YOUR_KEY"
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: Could not load .env file. Ensure GEMINI_API_KEY is set in your environment.")
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: GEMINI_API_KEY not found. Please set it in your environment or .env file.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 2. Initialize RAG components
	geminiService, err := NewGeminiService(ctx, apiKey)
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	// Mock Retriever for demonstration, replace with a real Vector DB client in production
	retriever := &MockRetriever{}

	// The GeminiService acts as both the Embedder and the Generator
	embedder := geminiService
	generator := geminiService

	// 3. Define the user's query
	userQuery := "What is the policy for submitting Paid Time Off?"
	log.Printf("User Query: %s\n", userQuery)

	// 4. Run the full RAG pipeline
	response, err := RunRAG(ctx, userQuery, retriever, embedder, generator)
	if err != nil {
		log.Fatalf("RAG pipeline failed: %v", err)
	}

	// 5. Output the result
	log.Println("--- RAG Response ---")
	log.Println(response)
}
