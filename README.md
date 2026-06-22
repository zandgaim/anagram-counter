# Anagram Counter

A concurrent Go application designed to process large text files and aggregate anagram groups. 

## Assumptions
* **Memory Management:** To maintain bounded memory when processing GBs of files, the application stores and outputs one representative "example word" and the total count for each anagram group, rather than keeping all string variants in RAM. Memory is bounded by the number of *unique* anagram signatures, rather than being strictly bounded overall.
* **Words & Punctuation:** A word is defined as alphanumeric characters separated by spaces or punctuation. Punctuation is used purely as a delimiter and is stripped. Words are case-insensitive. Unicode letters are supported intentionally via rune sorting.
* **Files:** UTF-8 encoded, immutable during execution, and can exceed available RAM. Processed chunk-by-chunk via buffered I/O. Unreadable files generate an error that is returned and reported, and processing terminates gracefully.
* **Concurrency:** Bounded worker pool defaulting to the host machine's logical CPU cores to prevent disk thrashing.
* **Output:** Results are formatted as JSON and sorted deterministically (by count descending, then alphabetically by the example word).

## Architecture
* **Domain (`internal/domain/anagram.go`):** Pure logic for anagram signature generation and state structures.
* **Operation (`internal/operation/reader.go`):** Concurrent disk I/O, token scanning, punctuation splitting, and worker pool management.
* **Facade (`internal/facade/orchestrator.go`):** Global state orchestration, progress reporting, and final JSON output formatting.

## Build and Run (Linux / Go 1.20+)

### Using Make (Recommended)
```bash
# Compile the binary
make build

# Run unit tests
make test

# Execute against a test directory
make run
