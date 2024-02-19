package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func downloadFile(url, filepath string, wg *sync.WaitGroup, errors chan<- error) {
	defer wg.Done()

	// Get the data from the URL
	resp, err := http.Get(url)
	if err != nil {
		errors <- err
		return
	}
	defer resp.Body.Close()

	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		errors <- err
		return
	}
	defer file.Close()

	// Write the data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		errors <- err
		return
	}

	fmt.Printf("Downloaded %s to %s\n", url, filepath)
}

func main() {
	// List of URLs to download
	urls := []string{
		"https://example.com/file1.txt",
		"https://example.com/file2.txt",
		"https://example.com/file3.txt",
		"https://example.com/file4.txt",
		"https://example.com/file5.txt",
	}

	// Output directory for downloaded files
	outputDir := "./downloads/"

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	var wg sync.WaitGroup
	errors := make(chan error)

	// Download files concurrently
	for _, url := range urls {
		wg.Add(1)
		go downloadFile(url, outputDir+extractFilename(url), &wg, errors)
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	// Check for any download errors
	for err := range errors {
		fmt.Printf("Error downloading: %s\n", err)
	}
}

func extractFilename(url string) string {
	// Extract filename from URL
	// Example: https://example.com/file.txt -> file.txt
	filename := url
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			filename = url[i+1:]
			break
		}
	}
	return filename
}
