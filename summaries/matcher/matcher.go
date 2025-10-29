package main

import (
	"context"
	"runtime"
	"sort"

	"github.com/philippgille/chromem-go"
	"k8s.io/klog/v2"
)

// PrepareCollection gets an existing chromem-go collection or creates a new one if it doesn't exist.
// It preserves any existing collection to avoid losing previously added documents.
//
// The function uses the configured embedding model (typically "nomic-embed-text")
// via Ollama for generating document embeddings. The collection is configured
// to use all available CPU cores for parallel processing.
//
// Parameters:
//   - database: The chromem-go database instance
//   - collectionName: Name of the collection to get or create
//
// Returns:
//   - *chromem.Collection: The collection ready for document operations
//   - error: Any error encountered during collection retrieval or creation
func PrepareCollection(database *chromem.DB, collectionName string) (*chromem.Collection, error) {
	embeddingFunc := chromem.NewEmbeddingFuncOllama(EmbeddingModel, "")

	collection, err := database.GetOrCreateCollection(collectionName, nil, embeddingFunc)
	if err != nil {
		return nil, err
	}

	// Check if this is a new collection or existing one by checking if it has documents
	// Note: We can't easily check this without querying, so we'll log based on whether
	// GetCollection would have returned it (which GetOrCreateCollection handles internally)
	klog.V(90).Infof("Prepared collection %q", collectionName)

	return collection, nil
}

// AddDocumentsToCollection adds a batch of documents to the specified collection
// and generates embeddings for them. The operation uses parallel processing
// with a number of workers equal to the available CPU cores.
//
// Parameters:
//   - collection: The chromem-go collection to add documents to
//   - docs: Slice of documents to add and embed
//
// Returns:
//   - error: Any error encountered during document addition or embedding generation
func AddDocumentsToCollection(collection *chromem.Collection, docs []chromem.Document) error {
	return collection.AddDocuments(context.TODO(), docs, runtime.NumCPU())
}

// FindSimilarDocuments queries the collection to find the most similar documents
// to the given document. It constructs a search query from the document's content
// and filters out the document itself from the results.
//
// The function filters out documents that have the same first path element as
// the document being compared by using a where clause that excludes documents
// with matching metadata.
//
// The function queries for DefaultTopSimilarCount results (6) to account for
// filtering out the document itself, then returns up to MaxSimilarResults (5)
// most similar documents sorted by similarity score in descending order.
//
// Parameters:
//   - collection: The chromem-go collection to query
//   - doc: The document to find similar documents for
//
// Returns:
//   - []SimilarityResult: Slice of similar documents with their similarity scores
//   - error: Any error encountered during the query operation
func FindSimilarDocuments(collection *chromem.Collection, doc chromem.Document) ([]SimilarityResult, error) {
	// Extract first path element and create where clause to exclude documents
	// with the same first path element
	firstPathElement := ExtractFirstPathElement(doc.ID)
	where := make(map[string]string)

	if firstPathElement != "" {
		// Use empty string as value to exclude documents where this metadata key
		// equals the first path element (since we set metadata[key] = key)
		where[firstPathElement] = ""
	}

	similarDocs, err := collection.Query(context.TODO(), doc.Content, DefaultTopSimilarCount, where, nil)
	if err != nil {
		return nil, err
	}

	// Filter out the document itself and limit to top 5
	similar := make([]SimilarityResult, 0, MaxSimilarResults)

	for _, similarDoc := range similarDocs {
		if similarDoc.ID != doc.ID {
			similar = append(similar, SimilarityResult{
				DocumentID: similarDoc.ID,
				Similarity: float64(similarDoc.Similarity),
			})
			if len(similar) >= MaxSimilarResults {
				break
			}
		}
	}

	return similar, nil
}

// ComputeSimilarityResults processes all documents in the collection and computes
// similarity results for each one. For each document, it finds the top 5 most
// similar documents (excluding itself) and aggregates the results into a flat list
// of similarity pairs, sorted by similarity in descending order.
//
// This function performs O(n) queries where n is the number of documents,
// making it suitable for batch processing but potentially slow for large datasets.
//
// Parameters:
//   - collection: The chromem-go collection containing the documents
//   - docs: Slice of all documents to process
//
// Returns:
//   - Results: Flat list of similarity pairs sorted by similarity in descending order
//   - error: Any error encountered during similarity computation
func ComputeSimilarityResults(collection *chromem.Collection, docs []chromem.Document) (Results, error) {
	pairs := make([]SimilarityPair, 0)

	for i, doc := range docs {
		klog.V(90).Infof("Processing document %d/%d: %s", i+1, len(docs), doc.ID)

		similar, err := FindSimilarDocuments(collection, doc)
		if err != nil {
			return nil, err
		}

		// Convert each similarity result into a pair
		for _, sim := range similar {
			pairs = append(pairs, SimilarityPair{
				Similarity: sim.Similarity,
				DocumentA:  doc.ID,
				DocumentB:  sim.DocumentID,
			})
		}
	}

	// Sort by similarity in descending order
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Similarity > pairs[j].Similarity
	})

	return Results(pairs), nil
}
