package operation

import (
    "bufio"
    "io"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "unicode"

    "github.com/zandgaim/anagram-counter/internal/domain"
)

type Result struct {
    Signatures   map[string]domain.GroupStats
    WordsScanned int
    Errors       error
}

// DiscoverFiles finds all .txt files in the given directory.
func DiscoverFiles(dir string) (<-chan string, <-chan error) {
    filesChan := make(chan string)
    errChan := make(chan error, 1)

    go func() {
        defer close(filesChan)
        defer close(errChan)

        err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
            if err != nil {
                return err
            }
            if !d.IsDir() && strings.EqualFold(filepath.Ext(path), ".txt") {
                filesChan <- path
            }
            return nil
        })

        if err != nil {
            errChan <- err
        }
    }()

    return filesChan, errChan
}

// ProcessFiles manages the worker pool.
func ProcessFiles(filePaths <-chan string, concurrency int, resultsChan chan<- Result) {
    var wg sync.WaitGroup

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for path := range filePaths {
                file, err := os.Open(path)
                if err != nil {
                    resultsChan <- Result{Errors: err}
                    continue
                }
                processStream(file, resultsChan)
                if err := file.Close(); err != nil {
                    resultsChan <- Result{Errors: err}
                }
            }
        }()
    }

    wg.Wait()
    close(resultsChan)
}

// processStream takes an io.Reader to make it testable without actual disk files.
func processStream(reader io.Reader, resultsChan chan<- Result) {
    scanner := bufio.NewScanner(reader)
    scanner.Buffer(make([]byte, 1024), 1024*1024)
    scanner.Split(bufio.ScanWords)

    localSigs := make(map[string]domain.GroupStats)
    wordsScanned := 0

    for scanner.Scan() {
        rawToken := scanner.Text()

        // Split the token by any non-alphanumeric character (handles punctuation attached to words)
        words := strings.FieldsFunc(rawToken, func(r rune) bool {
            return !unicode.IsLetter(r) && !unicode.IsNumber(r)
        })

        for _, w := range words {
            cleanWord := strings.ToLower(w)
            if cleanWord != "" {
                sig := domain.Signature(cleanWord)
                stats := localSigs[sig]
                
                // Save the first word we find as the representative example
                if stats.Count == 0 {
                    stats.ExampleWord = cleanWord
                }
                
                stats.Count++
                localSigs[sig] = stats
                wordsScanned++
            }
        }
    }

    if err := scanner.Err(); err != nil {
        resultsChan <- Result{Errors: err}
        return
    }

    resultsChan <- Result{
        Signatures:   localSigs,
        WordsScanned: wordsScanned,
    }
}