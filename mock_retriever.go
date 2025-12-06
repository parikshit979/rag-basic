package main

import (
	"context"
	"math/rand"
)

// MockRetriever simulates a Vector Database by returning a fixed context.
type MockRetriever struct{}

// Search implements the Retriever interface.
func (m *MockRetriever) Search(ctx context.Context, queryVector []float32, k int) ([]string, error) {
	// In a real implementation, this queryVector would be used to search.
	// Here, we ignore the vector and return a fixed context for demonstration.

	rand.Seed(rand.Int63()) // Simple seed for demonstration purposes

	if rand.Intn(10) < 2 { // Simulate retrieval failure 20% of the time
		return []string{}, nil
	}

	// Hardcoded context that the LLM will use to ground its answer
	context := []string{
		"The official company policy states that Paid Time Off (PTO) requests must be submitted through the internal HR portal at least 14 days in advance.",
		"The 'Employee Handbook 2024' specifies that for PTO requests exceeding 5 consecutive days, approval from a direct manager and HR is required.",
		"PTO balances can be checked anytime via the 'Benefits' tab on the employee dashboard.",
	}

	// Return the top k relevant documents (or all, in this mock example)
	return context[:min(k, len(context))], nil
}
