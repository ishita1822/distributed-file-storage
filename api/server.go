package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	file "distributed-file-storage/entity"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/files", file.GetFilesHandler).Methods("GET")
	r.HandleFunc("/files/{file_id}", file.GetFileHandler).Methods("GET")
	r.HandleFunc("/upload", file.UploadFileHandler).Methods("POST")
	r.HandleFunc("/files/{file_id}", file.UpdateFileHandler).Methods("PUT")
	r.HandleFunc("/files/{file_id}", file.DeleteFileHandler).Methods("DELETE")
	return r
}

func StartServer() {
	r := SetupRoutes()
	log.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", r)
}
