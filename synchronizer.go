package main

import "sync"

func synchronizeChunks(chunks []ChunkMeta, uploader Uploader, wg *sync.WaitGroup, generatedFileName string) error {
	// Create channels to communicate between goroutines
	chunkChan := make(chan ChunkMeta, len(chunks)) // Channel to send chunks to workers
	errChan := make(chan error, len(chunks))       // Channel to receive errors from workers

	// Iterate over the chunks slice and send each chunk to the chunk channel
	for _, chunk := range chunks {
		wg.Add(1)
		chunkChan <- chunk
	}

	close(chunkChan) // Close the chunk channel to signal that all chunks have been sent

	// Start multiple goroutines to process chunks in parallel
	for i := 0; i < len(chunks); i++ { // Number of parallel workers
		go func() {
			totalChunk := len(chunks)
			for chunk := range chunkChan { // Iterate over chunks received from the chunk channel
				defer wg.Done() // Decrease the WaitGroup counter when the goroutine finishes

				// Upload the chunk using the uploader interface
				err := uploader.UploadChunk(chunk, totalChunk, generatedFileName)
				if err != nil {
					errChan <- err // Send any errors to the error channel
					return
				}

			}
		}()
	}

	wg.Wait()      // Wait for all goroutines to finish processing chunks
	close(errChan) // Close the error channel after all errors have been received

	// Check for errors from the error channel
	for err := range errChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	return nil // Return nil if no errors occurred during synchronization
}
