package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func determinePathToClean() string {
	// Define the Downloads folder path
	downloadsPath := filepath.Join(os.Getenv("HOME"), "<Specify the folder full path>")
	return downloadsPath
}

func buildToDeletePath(foldername string) string {
	downloadsPath := determinePathToClean()
	// Define the to_delete folder path
	toDeletePath := filepath.Join(downloadsPath, foldername)

	// Create the to_delete folder if it doesn't exist
	if _, err := os.Stat(toDeletePath); os.IsNotExist(err) {
		if err := os.Mkdir(toDeletePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s Folder is created\n", toDeletePath)
	} else {
		fmt.Println("The folder already exists. Skipping folder creation.")
	}
	return toDeletePath
}

func moveFilesIntoToDelete(directory string, toDel string, maxAgeDays int) {
	// Get the current time
	now := time.Now()
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err1 error) error {
		// Skip the to_delete folder
		if path == toDel {
			return filepath.SkipDir
		}

		// Retrieve file info so we can retrieve the modTime
		fileInfo, err := os.Stat(path)
		if err != nil {
			fmt.Println("Error getting file info:", err)
		}

		// Get the modification file date
		modTime := fileInfo.ModTime()

		// Calculate the days from the last modification
		ageDays := int(now.Sub(modTime).Hours() / 24)

		if err != nil {
			return err
		}

		if ageDays > maxAgeDays {
			filename := filepath.Base(path)
			newPath := filepath.Join(toDel, filename)
			err := os.Rename(path, newPath)

			if err != nil {
				fmt.Printf("Error moving file %s: %s\n", filename, err)
			} else {
				fmt.Printf("Moved file %s to %s\n", filename, toDel)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", directory, err)
		return
	}
}

func main() {
	downloadsPath := determinePathToClean()
	toDeletePath := buildToDeletePath("to_delete")
	moveFilesIntoToDelete(downloadsPath, toDeletePath, 10)
}
