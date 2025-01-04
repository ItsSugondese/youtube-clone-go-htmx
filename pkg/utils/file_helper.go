package utils

import (
	"cloud.google.com/go/storage"
	"context"
	filepathconstants "youtube-clone/constants/file_path_constants"
	"youtube-clone/constants/file_type_constants"
	globaldto "youtube-clone/global/global_dto"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ClientUploader struct {
	Cl *storage.Client
	//ProjectID  string
	BucketName string
	UploadPath string
}

var Uploader *ClientUploader

// SaveFile saves the uploaded file to the specified directory and returns the URL of the saved file.
func SaveFile(file *multipart.FileHeader, module string, forBucket bool) globaldto.FileDetails {
	var uploadDir string

	if forBucket {
		uploadDir = filepathconstants.FilePathMappings[module].Path
	} else {
		uploadDir = filepath.Join(filepathconstants.UploadDir, filepathconstants.FilePathMappings[module].Path)
	}
	// Create the upload directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			panic("unable to create directory: " + err.Error())
		}
	}

	fileType := validateExtension(file.Filename)
	// Create a unique file name based on the current timestamp
	timestamp := time.Now().UnixNano()
	extension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", timestamp, extension)

	filePath := filepath.Join(uploadDir, newFileName)

	// SAVE the file to the specified directory
	if err := saveUploadedFile(file, filePath, forBucket); err != nil {
		panic("unable to save the file: " + err.Error())
	}

	// Return the URL of the saved file
	//fileURL := "localhost:3000/images/" + newFileName
	return globaldto.FileDetails{
		FilePath: filePath,
		Size:     file.Size,
		FileType: fileType,
	}
}

// saveUploadedFile is a helper function to save the uploaded file to the file system
func saveUploadedFile(file *multipart.FileHeader, filePath string, forBucket bool) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if forBucket {
		ctx := context.Background()

		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()
		wc := Uploader.Cl.Bucket(Uploader.BucketName).Object(filePath).NewWriter(ctx)
		defer wc.Close()
		if _, err := io.Copy(wc, src); err != nil {
			return fmt.Errorf("io.Copy: %v", err)
		}
		if err := wc.Close(); err != nil {
			return (fmt.Errorf("Writer.Close: %v", err))
		}
		return nil
	} else {

		out, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = out.ReadFrom(src)
		return err
	}
}

// fins and return the file from the path
func GetFileFromFilePath(filePath string, w http.ResponseWriter, fromBucket bool) {
	if filePath == "" {
		panic("File path is required")

	}

	fileName := filepath.Base(filePath)
	if fileName == "" {
		panic("Invalid file name")

	}

	if fromBucket {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()

		// Open the object from GCS
		reader, err := Uploader.Cl.Bucket(Uploader.BucketName).Object(filePath).NewReader(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to open GCS object: %v", err))
		}
		defer reader.Close()

		// Copy the GCS object to the response writer
		if _, err := io.Copy(w, reader); err != nil {
			panic(fmt.Errorf("io.Copy: %v", err))
		}
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			panic("Invalid file path")

		}
		defer file.Close()

		_, err = io.Copy(w, file)
		if err != nil {
			panic("Failed to write file to response")

		}
	}
}

// responsible for copying file from one path to another. will primiralrly be used to copy from temporary file to actual file path
func CopyFileToServer(filePath string, fileToPath string, toBucket bool) string {
	var overallToPath string

	fileName := filepath.Base(filePath)
	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")

	if toBucket {

		// Format the date as YYYY-MM-DD
		overallToPath = filepath.Join(filepathconstants.FilePathMappings[fileToPath].Location, date, fileName)

		// Copy the file
		err := CopyFileToGCS(filePath, overallToPath)
		if err != nil {
			panic("Failed to copy the file: " + err.Error())
		}

	} else {
		fileTo := filepath.Join(filepathconstants.UploadDir, filepathconstants.FilePathMappings[fileToPath].Location, date)

		overallToPath = filepath.Join(fileTo, fileName)

		// Create the directory if it doesn't exist
		err := os.MkdirAll(fileTo, os.ModePerm)
		if err != nil {
			panic("Failed to create directory: " + err.Error())
		}

		// Copy the file
		err = copyFile(filePath, overallToPath)
		if err != nil {
			panic("Failed to copy the file: " + err.Error())
		}
	}

	return overallToPath
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func CopyFileToGCS(filePath string, fileToPath string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	src := Uploader.Cl.Bucket(Uploader.BucketName).Object(filePath)
	dst := Uploader.Cl.Bucket(Uploader.BucketName).Object(fileToPath)

	// Open the source file for reading
	srcReader, err := src.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create source reader: %v", err)
	}
	defer srcReader.Close()

	// Open the destination file for writing
	dstWriter := dst.NewWriter(ctx)
	defer dstWriter.Close()

	// Copy the file
	if _, err := io.Copy(dstWriter, srcReader); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	return nil
}

// for validating whether the extension of the file is valid or not
func validateExtension(filename string) file_type_constants.FileType {
	extension := strings.ToUpper(filepath.Ext(filename))[1:] // get file extension without dot

	// Check if the extension is empty
	if extension == "" {
		panic("file has no extension")
	}

	var fileType file_type_constants.FileType

	if fileType, ok := file_type_constants.ImageType[extension]; ok {
		return fileType
	} else if fileType, ok := file_type_constants.DocumentType[extension]; ok {
		return fileType
	} else if fileType, ok := file_type_constants.PdfType[extension]; ok {
		return fileType
	} else if fileType, ok := file_type_constants.TxtType[extension]; ok {
		return fileType
	} else if fileType, ok := file_type_constants.ExcelType[extension]; ok {
		return fileType
	} else {
		panic("Not a valid extension")
	}

	// Prepare the result map
	return fileType
}
