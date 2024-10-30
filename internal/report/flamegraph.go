package main

import (
	"encoding/json"
	"os"
)

// FlameGraphTree represents a tree of test suites in the format required by the d3-flame-graph library.
type FlameGraphTree struct {
	Name     string            `json:"name"`
	Value    int               `json:"value"`
	Children []*FlameGraphTree `json:"children"`
}

// NewFromSuiteTree creates a new FlameGraphTree from a SuiteTree. The value of each node in the FlameGraphTree is the
// number of specs in the corresponding SuiteTree node.
func NewFromSuiteTree(tree *SuiteTree) *FlameGraphTree {
	root := &FlameGraphTree{
		Name:  tree.Name,
		Value: tree.Specs,
	}

	for _, child := range tree.Children {
		root.Children = append(root.Children, NewFromSuiteTree(child))
	}

	return root
}

// Save saves the FlameGraphTree as JSON to the given path.
func (tree *FlameGraphTree) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(tree)
}
