package entity

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"distributed-file-storage/model"
)

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := GetAllFiles()
	if err != nil {
		http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(files)
}

func GetFileHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	fileID, ok := vars["file_id"]
	if !ok {
		http.Error(w, "file_id is required", http.StatusBadRequest)
		return
	}

	files, err := GetFile(fileID)
	if err != nil {
		http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(files)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB file size limit
	if err != nil {
		log.Println("Error parsing multipart form:", err)
		http.Error(w, "File too large or invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve file from form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file from request:", err)
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	log.Println("Uploading file:", handler.Filename)

	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Println("Error creating upload directory:", err)
		http.Error(w, "Could not create upload directory", http.StatusInternalServerError)
		return
	}

	// Define file path
	filePath := uploadDir + handler.Filename

	// Open destination file
	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file on disk:", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the entire file without modifying line breaks
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println("Error writing file to disk:", err)
		http.Error(w, "Error writing file to disk", http.StatusInternalServerError)
		return
	}

	// Insert file metadata into the database
	uploadedFile := model.File{
		Name:      handler.Filename,
		Path:      filePath,
		CreatedAt: time.Now(),
	}

	err = InsertFile(uploadedFile)
	if err != nil {
		log.Println("Database insert error:", err)
		http.Error(w, "Failed to save file info in DB", http.StatusInternalServerError)
		return
	}

	log.Println("File uploaded successfully:", filePath)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "File uploaded successfully",
		"file_path": filePath,
	})
}

func UpdateFileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract file_id from the URL
	vars := mux.Vars(r)
	fileID, ok := vars["file_id"]
	if !ok {
		http.Error(w, "file_id is required", http.StatusBadRequest)
		return
	}

	// Decode JSON request body
	var req model.File
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close body after reading

	fileId, err := uuid.Parse(fileID)
	if err != nil {
		log.Fatalf("Invalid UUID: %v", err)
	}
	req.ID = fileId
	req.CreatedAt = time.Now()

	//validation if fileId exists in db
	if _, err = GetFile(fileID); err != nil {
		log.Printf("File with ID %s not found: %v", fileID, err) // Log the error for debugging
		http.Error(w, "File not found", http.StatusNotFound)     // Return 404
		return
	}

	// Call function to update the database
	err = updateFileRecord(fileID, req)
	if err != nil {
		http.Error(w, "Failed to update file", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File updated successfully"))
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract file_id from URL
	vars := mux.Vars(r)
	fileID, ok := vars["file_id"]
	if !ok {
		http.Error(w, "file_id is required", http.StatusBadRequest)
		return
	}

	err := deleteFileHandler(fileID)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusNoContent)
	fmt.Println("File deleted successfully!")
}
