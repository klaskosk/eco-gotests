package main

import (
	"encoding/json"
	"os"

	"github.com/philippgille/chromem-go"
	"k8s.io/klog/v2"
)

func init() {
	InitConfig()
}

func main() {
	klog.V(90).Infof("Running matcher tool and looking for docs in %q and storing the database in %q",
		GetDocDir(), GetDBDir())
	klog.V(90).Infof("Output will be saved to %q", GetOutputFile())
	klog.V(90).Infof("Tracking file: %q", GetTrackingFile())

	klog.V(90).Info("Loading document tracker")

	tracker, err := LoadTracker(GetTrackingFile())
	if err != nil {
		klog.Fatalf("Failed to load document tracker: %v", err)
	}

	klog.V(90).Infof("Loaded %d tracked documents", tracker.Count())

	klog.V(90).Info("Setting up chromem-go")

	database, err := chromem.NewPersistentDB(GetDBDir(), true)
	if err != nil {
		klog.Fatalf("Failed to create chromem-go database: %v", err)
	}

	klog.V(90).Info("Reading all documents from directory")

	docs, err := ReadSummariesIntoDocuments(GetDocDir())
	if err != nil {
		klog.Fatalf("Failed to read summaries into documents: %v", err)
	}

	if len(docs) == 0 {
		klog.Fatalf("No documents found in directory %q", GetDocDir())
	}

	klog.V(90).Infof("Found %d documents in directory", len(docs))

	// Filter out documents that are already tracked
	newDocs := tracker.FilterNewDocuments(docs)
	klog.V(90).Infof("Found %d new documents to add (skipping %d already tracked)", len(newDocs), len(docs)-len(newDocs))

	if len(newDocs) == 0 {
		klog.V(90).Info("No new documents to add, all documents are already in the database")
	} else {
		collection, err := PrepareCollection(database, CollectionName)
		if err != nil {
			klog.Fatalf("Failed to prepare collection: %v", err)
		}

		klog.V(90).Infof("Adding %d new documents to collection", len(newDocs))

		err = AddDocumentsToCollection(collection, newDocs)
		if err != nil {
			klog.Fatalf("Failed to add summaries to collection: %v", err)
		}

		// Update tracker with newly added documents
		newDocIDs := make([]string, len(newDocs))
		for i, doc := range newDocs {
			newDocIDs[i] = doc.ID
		}

		tracker.Add(newDocIDs...)

		klog.V(90).Info("Saving updated tracker")

		if err := tracker.Save(); err != nil {
			klog.Fatalf("Failed to save tracker: %v", err)
		}
	}

	klog.V(90).Info("Finding top 5 most similar documents for each document")

	// Use all documents (not just new ones) for similarity computation
	collection, err := PrepareCollection(database, CollectionName)
	if err != nil {
		klog.Fatalf("Failed to prepare collection: %v", err)
	}

	results, err := ComputeSimilarityResults(collection, docs)
	if err != nil {
		klog.Fatalf("Failed to compute similarity results: %v", err)
	}

	klog.V(90).Info("Saving results to JSON file")

	if err := SaveResultsToFile(results, GetOutputFile()); err != nil {
		klog.Fatalf("Failed to save results: %v", err)
	}

	klog.V(90).Infof("Successfully saved results to %q", GetOutputFile())
}

// SaveResultsToFile marshals the similarity results to JSON and writes them to a file.
// The JSON output is indented with 2 spaces for readability. The file is created
// with permissions 0644 (readable by all, writable by owner).
//
// Parameters:
//   - results: The similarity results to save (already sorted by similarity in descending order)
//   - filePath: The path where the JSON file should be written
//
// Returns:
//   - error: Any error encountered during JSON marshaling or file I/O
func SaveResultsToFile(results Results, filePath string) error {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}
