# Anagram Counter

A concurrent Go application designed to process large text files and aggregate anagram groups. Note that memory is bounded by the number of unique anagram signatures, rather than being strictly bounded overall.

## Assumptions
* **Words:** Alphanumeric only, punctuation is stripped, case-insensitive. Numbers are treated as part of words. Unicode letters are supported intentionally (note: we use sorting for anagram signatures which correctly handles full Unicode; for ASCII-only requirements, an array-based frequency count would be faster).
* **Files:** UTF-8 encoded, immutable during execution, and can exceed available RAM. Processed chunk-by-chunk via buffered I/O. Unreadable files generate an error that is returned and reported, and processing terminates gracefully.
* **Concurrency:** Bounded worker pool defaulting to the host machine's logical CPU cores to prevent disk thrashing.
* **Output:** Output will include the word signatures and counts, formatted deterministically (sorted).

## Architecture
* **Domain (`internal/domain`):** Pure logic for word sanitization and anagram signature generation.
* **Operation (`internal/operation`):** Concurrent disk I/O and worker pool.
* **Facade (`internal/facade`):** Global state orchestration and progress reporting.

## Build and Run (Linux / Go 1.20+)

### Using Make (Recommended)
```bash
# Compile the binary
make build

# Run unit tests
make test

# Execute against a test directory
make run

### Quality Checks
```bash
go test ./...
go vet ./...
gofmt -w .
go test -race ./...
```
