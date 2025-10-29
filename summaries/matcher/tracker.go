package main

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/philippgille/chromem-go"
	"k8s.io/klog/v2"
)

// DocumentTracker maintains a set of document IDs that have been added to the database.
// It uses a map for O(1) lookup performance when checking if a document exists.
type DocumentTracker struct {
	mu        sync.RWMutex
	documents map[string]bool
	filePath  string
}

// LoadTracker reads the tracking file from disk and returns a DocumentTracker instance.
// If the file doesn't exist, it returns an empty tracker. If the file exists but
// cannot be read or parsed, it returns an error.
//
// Parameters:
//   - filePath: The path to the JSON tracking file
//
// Returns:
//   - *DocumentTracker: A tracker instance loaded with existing document IDs
//   - error: Any error encountered during file reading or JSON parsing
func LoadTracker(filePath string) (*DocumentTracker, error) {
	tracker := &DocumentTracker{
		documents: make(map[string]bool),
		filePath:  filePath,
	}

	// If file doesn't exist, return empty tracker
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		klog.V(90).Infof("Tracking file %q does not exist, starting with empty tracker", filePath)

		return tracker, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// If file is empty, return empty tracker
	if len(data) == 0 {
		return tracker, nil
	}

	var documentIDs []string
	if err := json.Unmarshal(data, &documentIDs); err != nil {
		return nil, err
	}

	tracker.mu.Lock()

	for _, id := range documentIDs {
		tracker.documents[id] = true
	}

	tracker.mu.Unlock()

	klog.V(90).Infof("Loaded %d tracked documents from %q", len(documentIDs), filePath)

	return tracker, nil
}

// Contains checks if a document ID is already tracked.
//
// Parameters:
//   - documentID: The document ID to check
//
// Returns:
//   - bool: true if the document is tracked, false otherwise
func (dt *DocumentTracker) Contains(documentID string) bool {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	return dt.documents[documentID]
}

// Add adds one or more document IDs to the tracker.
//
// Parameters:
//   - documentIDs: One or more document IDs to add
func (dt *DocumentTracker) Add(documentIDs ...string) {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	for _, id := range documentIDs {
		dt.documents[id] = true
	}
}

// FilterNewDocuments filters a slice of documents, returning only those that
// are not already tracked. This is useful for determining which documents
// need to be added to the database.
//
// Parameters:
//   - docs: Slice of documents to filter
//
// Returns:
//   - []chromem.Document: Slice containing only documents not yet tracked
func (dt *DocumentTracker) FilterNewDocuments(docs []chromem.Document) []chromem.Document {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	newDocs := make([]chromem.Document, 0)

	for _, doc := range docs {
		if !dt.documents[doc.ID] {
			newDocs = append(newDocs, doc)
		}
	}

	return newDocs
}

// Save writes the current set of tracked document IDs to the tracking file.
// The file is written as a JSON array of document ID strings, indented with
// 2 spaces for readability. The file is created with permissions 0644.
//
// Returns:
//   - error: Any error encountered during JSON marshaling or file I/O
func (dt *DocumentTracker) Save() error {
	dt.mu.RLock()

	documentIDs := make([]string, 0, len(dt.documents))
	for id := range dt.documents {
		documentIDs = append(documentIDs, id)
	}

	dt.mu.RUnlock()

	jsonData, err := json.MarshalIndent(documentIDs, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(dt.filePath, jsonData, 0644); err != nil {
		return err
	}

	klog.V(90).Infof("Saved %d tracked documents to %q", len(documentIDs), dt.filePath)

	return nil
}

// Count returns the number of documents currently tracked.
//
// Returns:
//   - int: The number of tracked documents
func (dt *DocumentTracker) Count() int {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	return len(dt.documents)
}
