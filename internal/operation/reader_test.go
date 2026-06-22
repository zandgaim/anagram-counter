package operation

import (
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "testing"
)

// TestProcessStream verifies the core chunking and counting logic in memory
func TestProcessStream(t *testing.T) {
    // Note the punctuation inside "hello,world!" to test the new FieldsFunc split
    input := "listen silent \n hello,world! \n h3ll0!!!"
    reader := strings.NewReader(input)

    resultsChan := make(chan Result, 1)
    go processStream(reader, resultsChan)
    res := <-resultsChan

    if res.Errors != nil {
        t.Fatalf("Expected no errors, got: %v", res.Errors)
    }

    if res.WordsScanned != 5 {
        t.Errorf("Expected 5 words scanned, got %d", res.WordsScanned)
    }

    if stats, exists := res.Signatures["eilnst"]; !exists || stats.Count != 2 {
        t.Errorf("Expected count of 2 for 'eilnst', got %d", stats.Count)
    }

    if stats, exists := res.Signatures["ehllo"]; !exists || stats.Count != 1 {
        t.Errorf("Expected count of 1 for 'ehllo', got %d", stats.Count)
    }
}

func TestProcessStream_LongToken(t *testing.T) {
    longToken := strings.Repeat("a", 100*1024)
    reader := strings.NewReader(longToken)

    resultsChan := make(chan Result, 1)
    go processStream(reader, resultsChan)
    res := <-resultsChan

    if res.Errors != nil {
        t.Fatalf("Expected no errors, got: %v", res.Errors)
    }

    if res.WordsScanned != 1 {
        t.Errorf("Expected 1 word scanned, got %d", res.WordsScanned)
    }
}

func TestDiscoverFiles(t *testing.T) {
    tempDir := t.TempDir()

    txtPath := filepath.Join(tempDir, "test.txt")
    os.WriteFile(txtPath, []byte("dummy text"), 0644)

    TXTPath := filepath.Join(tempDir, "test2.TXT")
    os.WriteFile(TXTPath, []byte("dummy text 2"), 0644)

    csvPath := filepath.Join(tempDir, "ignore.csv")
    os.WriteFile(csvPath, []byte("dummy,csv"), 0644)

    filesChan, errChan := DiscoverFiles(tempDir)

    var files []string
    for f := range filesChan {
        files = append(files, f)
    }

    err := <-errChan
    if err != nil {
        t.Fatalf("DiscoverFiles failed: %v", err)
    }

    if len(files) != 2 {
        t.Fatalf("Expected exactly 2 files discovered, got %d", len(files))
    }
}

func TestProcessFiles_UnreadableFile(t *testing.T) {
    tempDir := t.TempDir()
    txtPath := filepath.Join(tempDir, "unreadable.txt")

    os.Remove(txtPath)

    filesChan := make(chan string, 1)
    filesChan <- txtPath
    close(filesChan)

    resultsChan := make(chan Result, 1)
    ProcessFiles(filesChan, 1, resultsChan)

    res := <-resultsChan
    if res.Errors == nil {
        t.Fatalf("Expected an error for unreadable file, got nil")
    }
}

func TestProcessFiles_Concurrency(t *testing.T) {
    tempDir := t.TempDir()
    numFiles := 100

    for i := 0; i < numFiles; i++ {
        txtPath := filepath.Join(tempDir, "test_"+strconv.Itoa(i)+".txt")
        os.WriteFile(txtPath, []byte("listen silent"), 0644)
    }

    filesChan, errChan := DiscoverFiles(tempDir)
    resultsChan := make(chan Result, numFiles)

    ProcessFiles(filesChan, 4, resultsChan)

    err := <-errChan
    if err != nil {
        t.Fatalf("DiscoverFiles failed: %v", err)
    }

    filesProcessed := 0
    totalWords := 0

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        for res := range resultsChan {
            if res.Errors != nil {
                t.Errorf("Unexpected error: %v", res.Errors)
            } else {
                filesProcessed++
                totalWords += res.WordsScanned
            }
        }
    }()

    wg.Wait()

    if filesProcessed != numFiles {
        t.Errorf("Expected %d files processed, got %d", numFiles, filesProcessed)
    }
    if totalWords != numFiles*2 {
        t.Errorf("Expected %d words, got %d", numFiles*2, totalWords)
    }
}