package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type FileEntry struct {
	Index int
	File  string
}

func ReadFromFileAndAppendToSlice(filename string) []string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("Error opening file:", err))
	}
	defer file.Close()

	var entries []FileEntry

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into index and file parts
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue // Skip invalid lines
		}

		// Parse the index
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error parsing index:", err)
			continue
		}

		// Extract the file name
		fileName := strings.Trim(parts[2], "'")

		// Append to the entries slice
		entries = append(entries, FileEntry{
			Index: index,
			File:  fileName,
		})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("Error reading file:", err))
	}

	// Sort the entries by index
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Index < entries[j].Index
	})

	// Print or use the sorted entries
	for _, entry := range entries {
		fmt.Printf("Index: %d, File: %s\n", entry.Index, entry.File)
	}

	// If you need just the file names in a slice
	var fileNames []string
	for _, entry := range entries {
		fileNames = append(fileNames, entry.File)
	}
	fmt.Println("Sorted File Names:", fileNames)

	return fileNames
}

func ReadFromFileAndReturnFileEntrySlice(filename string) []FileEntry {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return []FileEntry{}
	}
	defer file.Close()

	var entries []FileEntry

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into index and file parts
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue // Skip invalid lines
		}

		// Parse the index
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error parsing index:", err)
			continue
		}

		// Extract the file name
		fileName := strings.Trim(parts[2], "'")

		// Append to the entries slice
		entries = append(entries, FileEntry{
			Index: index,
			File:  fileName,
		})
	}

	return entries
}

func CreateTempFile(fileEntries []string) string {
	// Create a temporary text file
	tempFile, err := ioutil.TempFile("", "ffmpeg_concat_*.txt")
	if err != nil {
		panic(fmt.Errorf("Error creating temporary file:", err))
	}
	defer os.Remove(tempFile.Name()) // Clean up after execution

	// Write slice data to the temporary file
	for _, entry := range fileEntries {
		_, err := tempFile.WriteString(entry + "\n")
		if err != nil {
			panic(fmt.Errorf("Error writing to temporary file:", err))
		}
	}

	// Close the temp file to flush changes
	tempFile.Close()

	return tempFile.Name()
}
