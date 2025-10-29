package main

import (
	"flag"

	"k8s.io/klog/v2"
)

const (
	// EmbeddingModel specifies the Ollama model used for generating document embeddings.
	// This model is used to convert text documents into vector representations for similarity matching.
	EmbeddingModel = "qwen3-embedding:0.6b"

	// CollectionName is the name of the chromem-go collection used to store documents.
	CollectionName = "summaries"

	// DefaultTopSimilarCount is the number of most similar documents to return for each query.
	// We query for 6 to account for filtering out the document itself, leaving 5 results.
	DefaultTopSimilarCount = 6

	// MaxSimilarResults is the maximum number of similar documents returned per document.
	MaxSimilarResults = 5
)

var (
	dbDir        *string
	docDir       *string
	outputFile   *string
	trackingFile *string
)

// InitConfig initializes command-line flags and logging configuration.
// This function must be called before using any configuration variables.
// It sets up flags for database directory, document directory, and output file path,
// and configures klog to output to stderr with verbose logging enabled.
func InitConfig() {
	dbDir = flag.String("db-dir", "./db", "Directory for the chromem-go database")
	docDir = flag.String("doc-dir", "./documents", "Directory containing summary documents")
	outputFile = flag.String("output", "./similarity_results.json", "Output JSON file path for similarity results")
	trackingFile = flag.String(
		"tracking-file", "./document_tracking.json", "JSON file path for tracking documents already in database")

	klog.InitFlags(nil)
	flag.Parse()

	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("v", "90")
}

// GetDBDir returns the configured database directory path.
func GetDBDir() string {
	return *dbDir
}

// GetDocDir returns the configured document directory path.
func GetDocDir() string {
	return *docDir
}

// GetOutputFile returns the configured output file path.
func GetOutputFile() string {
	return *outputFile
}

// GetTrackingFile returns the configured tracking file path.
func GetTrackingFile() string {
	return *trackingFile
}
