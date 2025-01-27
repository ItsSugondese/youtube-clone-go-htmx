package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
)

func (c *DefaultFileChunker) ChunkVideo(inputFile string, outputDir string, chunkDuration int, generatedFileName string) ([]ChunkMeta, error) {
	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %v", err)
	}

	// FFmpeg command to split the video
	outputPattern := filepath.Join(outputDir, generatedFileName+"_chunk_%03d.mp4")
	// Execute FFmpeg command
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c", "copy", // Copy codec (no re-encoding)
		"-map", "0", // Include all streams
		"-f", "segment", // Use segment format
		"-segment_time", strconv.Itoa(chunkDuration), // Chunk duration in seconds
		"-reset_timestamps", "1", // Reset timestamps for each segment
		outputPattern, // Output file pattern
	)

	// Run the FFmpeg command
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to chunk video: %v", err)
	}

	// Generate metadata for each	 chunk
	var chunks []ChunkMeta
	files, err := filepath.Glob(filepath.Join(outputDir, generatedFileName+"_chunk_*.mp4"))
	if err != nil {
		return nil, fmt.Errorf("failed to list chunk files: %v", err)
	}

	for index, file := range files {
		md5Hash, err := calculateMD5(file)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate MD5 for file %s: %v", file, err)
		}

		chunks = append(chunks, ChunkMeta{
			FileName: file,
			MD5Hash:  md5Hash,
			Index:    index,
		})
	}

	return chunks, nil
}

// ChunkFile splits a file into smaller chunks and returns metadata for each chunk.
// It reads the file sequentially and chunks it based on the specified chunk size.
func (c *DefaultFileChunker) ChunkFile(filePath string) ([]ChunkMeta, error) {
	var chunks []ChunkMeta // Store metadata for each chunk

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a buffer to hold the chunk data
	buffer := make([]byte, c.chunkSize)
	index := 0 // Initialize chunk index

	// Loop until EOF is reached
	for {
		// Read chunkSize bytes from the file into the buffer
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if bytesRead == 0 {
			break // If bytesRead is 0, it means EOF is reached
		}

		// Generate a unique hash for the chunk data
		hash := md5.Sum(buffer[:bytesRead])
		hashString := hex.EncodeToString(hash[:])

		// Construct the chunk file name
		chunkFileName := fmt.Sprintf("%s.chunk.%d", filePath, index)

		// Create a new chunk file and write the buffer data to it
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			return nil, err
		}
		_, err = chunkFile.Write(buffer[:bytesRead])
		if err != nil {
			return nil, err
		}

		metaData := ChunkMeta{FileName: chunkFileName, MD5Hash: hashString, Index: index}
		// Append metadata for the chunk to the chunks slice
		chunks = append(chunks, metaData)

		// Close the chunk file
		chunkFile.Close()

		// Move to the next chunk
		index++
	}

	return chunks, nil
}

// ChunklargeFile splits a large file into smaller chunks in parallel and returns metadata for each chunk.
// It divides the file into chunks and processes them concurrently using multiple goroutines.
func (c *DefaultFileChunker) ChunklargeFile(filePath string) ([]ChunkMeta, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var chunks []ChunkMeta // Store metadata for each chunk

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file information to determine the number of chunks
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	numChunks := int(fileInfo.Size() / int64(c.chunkSize))
	if fileInfo.Size()%int64(c.chunkSize) != 0 {
		numChunks++
	}

	// Create channels to communicate between goroutines
	chunkChan := make(chan ChunkMeta, numChunks)
	errChan := make(chan error, numChunks)
	indexChan := make(chan int, numChunks)

	// Populate the index channel with chunk indices
	for i := 0; i < numChunks; i++ {
		indexChan <- i
	}
	close(indexChan)

	// Start multiple goroutines to process chunks in parallel
	for i := 0; i < 4; i++ { // Number of parallel workers
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range indexChan {
				// Calculate the offset for the current chunk
				offset := int64(index) * int64(c.chunkSize)
				buffer := make([]byte, c.chunkSize) // Create a buffer for chunk data

				// Seek to the appropriate position in the file
				file.Seek(offset, 0)

				// Read chunkSize bytes from the file into the buffer
				bytesRead, err := file.Read(buffer)
				if err != nil && err != io.EOF {
					errChan <- err
					return
				}

				// If bytesRead is 0, it means EOF is reached
				if bytesRead > 0 {
					// Generate a unique hash for the chunk data
					hash := md5.Sum(buffer[:bytesRead])
					hashString := hex.EncodeToString(hash[:])

					// Construct the chunk file name
					chunkFileName := fmt.Sprintf("%s.chunk.%d", filePath, index)

					// Create a new chunk file and write the buffer data to it
					chunkFile, err := os.Create(chunkFileName)
					if err != nil {
						errChan <- err
						return
					}
					_, err = chunkFile.Write(buffer[:bytesRead])
					if err != nil {
						errChan <- err
						return
					}

					// Append metadata for the chunk to the chunks slice
					chunk := ChunkMeta{
						FileName: chunkFileName,
						MD5Hash:  hashString,
						Index:    index,
					}
					mu.Lock()
					chunks = append(chunks, chunk)
					mu.Unlock()

					// Close the chunk file
					chunkFile.Close()

					// Send the processed chunk to the chunk channel
					chunkChan <- chunk
				}
			}
		}()
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(chunkChan)
		close(errChan)
	}()

	// Check for errors from goroutines
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return chunks, nil
}

func calculateMD5(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a new MD5 hash
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to compute MD5 hash: %v", err)
	}

	// Return the MD5 hash as a hex string
	return hex.EncodeToString(hash.Sum(nil)), nil
}
