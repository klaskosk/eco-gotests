package main

import (
	_ "embed"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var (
	//go:embed tree_template.html
	treeTemplateFile string

	//go:embed report_template.html
	reportTemplateFile string
)

var (
	funcMap = template.FuncMap{
		"cleanPath": cleanPath,
	}
	treeTemplate   = template.Must(template.New("tree_template.html").Funcs(funcMap).Parse(treeTemplateFile))
	reportTemplate = template.Must(template.New("report_template.html").Parse(reportTemplateFile))
)

type TreeTemplateConfig struct {
	Tree       *SuiteTree
	Generated  time.Time
	Branch     string
	ActionURL  template.URL
	RepoURL    template.URL
	TimeFormat string
}

func TemplateTree(config TreeTemplateConfig, outputFileName string) error {
	return executeTemplateAndSave(treeTemplate, config, outputFileName)
}

type ReportTemplateConfig struct {
	BranchReports []BranchReportConfig
	Generated     time.Time
	ActionURL     template.URL
	RepoURL       template.URL
	TimeFormat    string
}

type BranchReportConfig struct {
	Name          string
	ReportFile    string
	Revision      string
	ShortRevision string
}

func TemplateReport(config ReportTemplateConfig, outputFileName string) error {
	slices.SortFunc(config.BranchReports, func(a, b BranchReportConfig) int {
		return strings.Compare(a.Name, b.Name)
	})

	return executeTemplateAndSave(reportTemplate, config, outputFileName)
}

// executeTemplateAndSave creates a file at outputFileName before executing tmpl with data provided by config. If
// outputFileName already exists, then it is truncated.
func executeTemplateAndSave(tmpl *template.Template, config any, outputFileName string) error {
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	err = tmpl.Execute(outputFile, config)
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
