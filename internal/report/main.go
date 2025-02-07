/*
Report is a tool to generate a report of the test suites in a Ginkgo test suite. It will print a tree of the test suites
and the number of specs in each suite. It will also generate a d3-flame-graph compatible JSON file.

Usage:

	report [flags]

The flags are:

	-h, -help
		Print this help message

	-b, -branch string
		Branch to clone and run tests for. Leave blank to use local repo

	-v int
		Log level verbosity for glog. Use 100 for logging all messages or leave blank for none
*/
package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
)

var (
	help   bool
	branch string
)

//nolint:gochecknoinits // This is a main package so init is fine
func init() {
	const (
		helpUsage   = "Print this help message"
		branchUsage = "Branch to clone and run tests for. Leave blank to use local repo"

		defaultHelp   = false
		defaultBranch = ""
	)

	flag.BoolVar(&help, "help", defaultHelp, helpUsage)
	flag.BoolVar(&help, "h", defaultHelp, helpUsage+" (shorthand)")

	flag.StringVar(&branch, "branch", defaultBranch, branchUsage)
	flag.StringVar(&branch, "b", defaultBranch, branchUsage+" (shorthand)")
}

func main() {
	// Also send glog messages to stderr
	_ = flag.Lookup("logtostderr").Value.Set("true")

	flag.Parse()

	if help {
		flag.Usage()

		return
	}

	cache, err := NewCache()
	if err != nil {
		glog.Errorf("Failed to create cache: %v", err)

		os.Exit(1)
	}

	var tree *SuiteTree

	if branch != "" {
		tree, err = getFromCacheOrClone(cache, "https://github.com/openshift-kni/eco-gotests", branch)
	} else {
		tree, err = getFromCacheOrDryRun(cache, ".")
	}

	if err != nil {
		glog.Errorf("Failed to get or create SuiteTree from cache: %v", err)

		os.Exit(1)
	}

	err = cache.Save()
	if err != nil {
		glog.Errorf("Failed to save cache: %v", err)

		os.Exit(1)
	}

	tree = tree.TrimRoot()
	tree.Sort(true)
	fmt.Print(tree)

	err = templateTree(templateValues{
		Tree:       tree,
		Generated:  time.Now(),
		Branch:     branch,
		ActionURL:  "/",
		RepoURL:    "https://github.com/openshift-kni/eco-gotests",
		TimeFormat: time.RFC3339,
	})
	if err != nil {
		glog.Errorf("Failed to template SuiteTree: %v", err)

		os.Exit(1)
	}

	flameGraph := NewFromSuiteTree(tree)
	err = flameGraph.Save("internal/report/data.json")

	if err != nil {
		glog.Errorf("Failed to save FlameGraphTree to data.json: %v", err)

		os.Exit(1)
	}
}

func getFromCacheOrClone(cache *Cache, repo, branch string) (*SuiteTree, error) {
	tree, err := cache.GetRemote(repo, branch)
	if err == nil {
		return tree, nil
	}

	repoPath, err := CloneRepo("/tmp", "https://github.com/openshift-kni/eco-gotests", branch)
	if err != nil {
		glog.Errorf("Failed to clone repo: %v", err)

		os.Exit(1)
	}

	tree, err = getFromCacheOrDryRun(cache, repoPath)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func getFromCacheOrDryRun(cache *Cache, repoPath string) (*SuiteTree, error) {
	tree, err := cache.GetOrCreate(repoPath, func() (*SuiteTree, error) {
		reportPath, err := DryRun(repoPath)
		if err != nil {
			glog.Errorf("Failed to run eco-gotests dry-run: %v", err)

			return nil, err
		}

		tree, err := NewFromFile(reportPath)
		if err != nil {
			glog.Errorf("Failed to create SuiteTree from report.json: %v", err)

			return nil, err
		}

		_ = os.Remove(reportPath)

		return tree, nil
	})
	if err != nil {
		glog.Errorf("Failed to get or create SuiteTree from cache: %v", err)

		return nil, err
	}

	return tree, nil
}

type templateValues struct {
	Tree       *SuiteTree
	Generated  time.Time
	Branch     string
	ActionURL  template.URL
	RepoURL    template.URL
	TimeFormat string
}

func templateTree(values templateValues) error {
	tmpl, err := template.New("template.html").
		Funcs(template.FuncMap{"cleanPath": cleanPath}).
		ParseFiles("internal/report/template.html")
	if err != nil {
		return err
	}

	outputFile, err := os.Create("internal/report/templated.html")
	if err != nil {
		return err
	}

	defer outputFile.Close()

	err = tmpl.Execute(outputFile, values)
	if err != nil {
		return err
	}

	return nil
}

// cleanPath cleans the provided path of anything preceding the eco-gotests directory.
func cleanPath(path string) string {
	pathElements := strings.Split(path, string(os.PathSeparator))
	for i, element := range pathElements {
		if element == "eco-gotests" {
			return filepath.Join(pathElements[i:]...)
		}
	}

	return path
}
