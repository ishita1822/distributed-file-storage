package entity

import (
	"database/sql"
	"fmt"
	"log"

	db "distributed-file-storage/db"
	"distributed-file-storage/model"
)

// var DB *sql.DB

// InsertFile saves file metadata to the database
func InsertFile(file model.File) error {

	query := `INSERT INTO files (name, path, created_at) VALUES ($1, $2, $3)`

	_, err := db.DB.Exec(query, file.Name, file.Path, file.CreatedAt)
	if err != nil {
		log.Println("Error inserting file:", err)
		return err
	}
	return nil
}

// GetAllFiles retrieves all files from the database
func GetAllFiles() ([]model.File, error) {
	rows, err := db.DB.Query("SELECT id, name, path, created_at FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []model.File
	for rows.Next() {
		var f model.File
		if err := rows.Scan(&f.ID, &f.Name, &f.Path, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

// GetFile retrieves a file record by fileID
func GetFile(fileID string) (model.File, error) {
	var file model.File

	// Prepare SQL query to fetch the file
	query := "SELECT id, name, path, created_at FROM files WHERE id = $1"

	// Execute the query and scan the result into the file struct
	err := db.DB.QueryRow(query, fileID).Scan(&file.ID, &file.Name, &file.Path, &file.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("File not found:", fileID)
			return model.File{}, err // Return empty file with no error to indicate "not found"
		}
		log.Println("Error fetching file from db:", err)
		return model.File{}, err
	}

	return file, nil
}
func updateFileRecord(fileID string, newContent model.File) error {
	// Prepare the UPDATE SQL query
	query := "UPDATE files SET name = $1, path = $2 WHERE id = $3"

	fmt.Println("here")
	// Execute the query
	_, err := db.DB.Exec(query, newContent.Name, newContent.Path, fileID)
	if err != nil {
		log.Println("Error updating file content in db:", err)
		return err
	}

	return nil
}

func deleteFileHandler(fileID string) error {
	// Prepare DELETE SQL query
	query := "DELETE FROM files WHERE id = $1"
	_, err := db.DB.Exec(query, fileID)
	if err != nil {
		log.Println("Error deleting file from db: ", err)
		return err
	}

	return err
}
