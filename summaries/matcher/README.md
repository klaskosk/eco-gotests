# matcher

This is a tool that uses the markdown files in the summaries directory to find similar or duplicate test cases.

## Prerequisites

This tool assumes that there is an ollama instance running on the local machine and available at `http://localhost:11434`. The embedding model used is `qwen3-embedding:0.6b` and it should be pulled before running the tool.

```bash
ollama pull qwen3-embedding:0.6b
```

## Running the tool

```bash
go build -o matcher.out ./matcher && ./matcher.out
```

## Output

The tool will output the similar or duplicate test cases in the `similarity_results.json` file and save the database in the `db` directory. This similarity results file is a list of similarity pairs between the test cases, in descending order of similarity. These similarity pairs are only between test cases that have different first path elements. This way, cnf/ran test cases get compared to system-tests test cases, but not to other cnf/ran test cases.

## Document Tracking

The tool maintains a tracking file that lists all documents already added to the database. This allows the tool to:

- Skip documents that have already been processed
- Only add new documents to the database on subsequent runs
- Preserve the existing database instead of rebuilding it from scratch

The tracking file is a JSON array of document IDs. If the file doesn't exist, it will be created automatically.

## Configuration

The tool can be configured with the following flags:

- `-db-dir`: The directory to save the database (default: `./db`)
- `-doc-dir`: The directory to read the documents from (default: `./documents`)
- `-output`: The file to save the similarity results to (default: `./similarity_results.json`)
- `-tracking-file`: The JSON file path for tracking documents already in database (default: `./document_tracking.json`)
