package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	logFileDir   = "./logfiles"
	outputDir    = "./filteredlogs"
	infoFile     = "info.log"
	warningFile  = "warning.log"
	errorFile    = "error.log"
	criticalFile = "critical.log"
)

var logLevels = [...]string{"info", "warning", "error", "critical"}

func processLogFiles() {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	logFiles, err := ioutil.ReadDir(logFileDir)
	if err != nil {
		log.Fatalf("Error reading log file directory: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(logFiles))

	for _, file := range logFiles {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			filePath := filepath.Join(logFileDir, file.Name())
			go processLogFile(filePath, &wg)
		}
	}

	wg.Wait()
	fmt.Println("Log files processed successfully.")
}

func processLogFile(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var fileWriters []struct {
		level string
		file  *os.File
	}

	for _, level := range logLevels {
		outputFilePath := filepath.Join(outputDir, level+".log")
		outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening output file %s: %v\n", outputFilePath, err)
			return
		}
		defer outputFile.Close()

		fileWriters = append(fileWriters, struct {
			level string
			file  *os.File
		}{level, outputFile})
	}

	for scanner.Scan() {
		line := scanner.Text()
		for _, fw := range fileWriters {
			if strings.Contains(line, "["+fw.level+"]") {
				_, err := fmt.Fprintln(fw.file, line)
				if err != nil {
					log.Printf("Error writing to output file: %v\n", err)
					return
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning file %s: %v\n", filePath, err)
		return
	}
}

func main() {
	processLogFiles()
}
