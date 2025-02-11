package main

import (
	"crypto/sha256"
	"embed"
	"fmt"
)

//go:embed *.go
var programSourceCode embed.FS

var sourceCodeSum string = getSourceSum()

func getSourceSum() string {
	summer := sha256.New()

	dirEntries, err := programSourceCode.ReadDir(".")
	if err != nil {
		return ""
	}

	for _, dirEntry := range dirEntries {
		contents, err := programSourceCode.ReadFile(dirEntry.Name())
		if err != nil {
			return ""
		}

		_, err = summer.Write(contents)
		if err != nil {
			return ""
		}
	}

	return fmt.Sprintf("%x", summer.Sum(nil))
}
