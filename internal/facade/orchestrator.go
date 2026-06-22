package facade

import (
    "encoding/json"
    "fmt"
    "runtime"
    "sort"
    "time"

    "github.com/zandgaim/anagram-counter/internal/domain"
    "github.com/zandgaim/anagram-counter/internal/operation"
)

func Run(directory string, concurrency int) error {
    fmt.Printf("Discovering files in %s...\n", directory)
    filesChan, errChan := operation.DiscoverFiles(directory)

    resultsChan := make(chan operation.Result)
    go operation.ProcessFiles(filesChan, concurrency, resultsChan)

    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    globalState := make(map[string]domain.GroupStats)
    totalWords := 0
    filesProcessed := 0
    errorsCount := 0

    startTime := time.Now()

    for {
        select {
        case err, ok := <-errChan:
            if !ok {
                errChan = nil // Prevents CPU busy-spin after channel closes
                continue
            }
            if err != nil {
                return fmt.Errorf("error discovering files: %w", err)
            }
        case res, ok := <-resultsChan:
            if !ok {
                printFinalReport(globalState, filesProcessed, totalWords, errorsCount, startTime)
                return nil
            }
            if res.Errors != nil {
                errorsCount++
            } else {
                filesProcessed++
                totalWords += res.WordsScanned
                for sig, stats := range res.Signatures {
                    globalStats := globalState[sig]
                    if globalStats.Count == 0 {
                        globalStats.ExampleWord = stats.ExampleWord
                    }
                    globalStats.Count += stats.Count
                    globalState[sig] = globalStats
                }
            }

        case <-ticker.C:
            fmt.Printf("[Progress] Time: %v | Files: %d | Words: %d | Errors: %d | Groups: %d\n",
                time.Since(startTime).Round(time.Second), filesProcessed, totalWords, errorsCount, len(globalState))
        }
    }
}

func printFinalReport(state map[string]domain.GroupStats, files, words, errors int, start time.Time) {
    fmt.Println("\n================ FINAL REPORT ================")
    fmt.Printf("Total Execution Time: %v\n", time.Since(start))
    fmt.Printf("Files Processed:      %d\n", files)
    fmt.Printf("Words Scanned:        %d\n", words)
    fmt.Printf("Errors Encountered:   %d\n", errors)
    fmt.Printf("Total Unique Words:   %d\n", len(state))

    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    // Using Sys instead of TotalAlloc to accurately reflect OS memory footprint, not cumulative allocations
    fmt.Printf("Total Memory Allocated (Sys): %v MB\n", m.Sys/1024/1024)
    fmt.Println("==============================================")

    type OutputGroup struct {
        ExampleWord string `json:"exampleWord"`
        Count       int    `json:"count"`
    }

    var groups []OutputGroup
    for _, stats := range state {
        if stats.Count > 1 { // Only record actual anagram groups (more than 1 occurrence)
            groups = append(groups, OutputGroup{ExampleWord: stats.ExampleWord, Count: stats.Count})
        }
    }

    // Sort by count descending, then alphabetically by the example word
    sort.Slice(groups, func(i, j int) bool {
        if groups[i].Count == groups[j].Count {
            return groups[i].ExampleWord < groups[j].ExampleWord
        }
        return groups[i].Count > groups[j].Count
    })

    output := struct {
        TotalWords    int           `json:"totalWords"`
        AnagramGroups []OutputGroup `json:"anagramGroups"`
    }{
        TotalWords:    words,
        AnagramGroups: groups,
    }

    out, _ := json.MarshalIndent(output, "", "  ")
    fmt.Println(string(out))
}