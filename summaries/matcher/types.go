package main

// SimilarityResult represents a single document similarity match with its score.
// The similarity score is a float64 value typically ranging from 0.0 to 1.0,
// where higher values indicate greater similarity.
type SimilarityResult struct {
	DocumentID string  `json:"document_id"`
	Similarity float64 `json:"similarity"`
}

// DocumentResults contains similarity results for a single document.
// It includes the document ID and a list of similar documents (excluding itself),
// sorted by similarity score in descending order.
type DocumentResults struct {
	DocumentID string             `json:"document_id"`
	Similar    []SimilarityResult `json:"similar"`
}

// SimilarityPair represents a pair of documents with their similarity score.
// This is used for the flat output format, sorted by similarity in descending order.
type SimilarityPair struct {
	Similarity float64 `json:"similarity"`
	DocumentA  string  `json:"document_a"`
	DocumentB  string  `json:"document_b"`
}

// Results is a list of similarity pairs, sorted by similarity in descending order.
type Results []SimilarityPair
