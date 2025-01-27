package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

func watchFile(filePath string, changeChan chan bool) {
	// Create a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err) // Terminate the program if an error occurs while creating the watcher
	}
	defer watcher.Close() // Close the watcher when the function exits

	// Add the specified file to the watcher's list of watched files
	err = watcher.Add(filePath)
	if err != nil {
		log.Fatal(err) // Terminate the program if an error occurs while adding the file to the watcher
	}

	// Infinite loop to continuously monitor events from the watcher
	for {
		select {
		case event, ok := <-watcher.Events:
			// Check if the events channel is closed
			if !ok {
				return // Exit the function if the channel is closed
			}
			// Check if the event corresponds to a write operation on the file
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("Modified file:", event.Name) // Log the name of the modified file
				changeChan <- true                        // Send a signal to the change channel indicating file modification
			}
		case err, ok := <-watcher.Errors:
			// Check if the errors channel is closed
			if !ok {
				return // Exit the function if the channel is closed
			}
			log.Println("Error:", err) // Log any errors that occur during file watching
		}
	}
}
