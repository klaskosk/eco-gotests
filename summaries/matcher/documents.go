package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/philippgille/chromem-go"
)

// ExtractFirstPathElement extracts the first underscore-delimited section
// from a document name. For example, "cnf_ran_ptp" returns "cnf".
// If there are no underscores, it returns the entire name.
func ExtractFirstPathElement(docName string) string {
	parts := strings.Split(docName, "_")
	if len(parts) > 0 {
		return parts[0]
	}

	return docName
}

// ReadSummariesIntoDocuments reads all regular files from the specified directory
// and converts them into chromem.Document objects. Each file's content is prefixed
// with "search_document: " to optimize it for semantic search queries.
//
// The function skips non-regular files (directories, symlinks, etc.) and uses
// the filename as the document ID. If any file cannot be read, the function
// returns an error immediately.
//
// Parameters:
//   - dir: The directory path containing the summary files to read
//
// Returns:
//   - []chromem.Document: A slice of document objects ready for embedding
//   - error: Any error encountered during directory reading or file I/O
func ReadSummariesIntoDocuments(dir string) ([]chromem.Document, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	docs := make([]chromem.Document, 0, len(files))

	for _, file := range files {
		if !file.Type().IsRegular() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		// Extract first path element and add to metadata
		firstPathElement := ExtractFirstPathElement(file.Name())
		metadata := make(map[string]string)
		metadata[firstPathElement] = firstPathElement

		docs = append(docs, chromem.Document{
			ID:       file.Name(),
			Content:  string(content),
			Metadata: metadata,
		})
	}

	return docs, nil
}
