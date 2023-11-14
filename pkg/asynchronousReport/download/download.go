package download

import (
	utils "FlexiProxyHub/internal/utils/generic"
	"FlexiProxyHub/pkg/database"
	"net/http"

	"go.uber.org/zap"
)

var log *zap.Logger
var DB *database.Database

func SetLog(logger *zap.Logger) {
	log = logger
}

// Set headers to dump file as download
func setHeaders(w http.ResponseWriter, filename string) {
	contentType, contentDisposition := database.GetContentTypeAndContentDisposition(filename)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", contentDisposition)
}

// Download handler and download file from local to client browser using cookie name filename
func Download(w http.ResponseWriter, r *http.Request) {
	filename, err := utils.GetCookie(r, "filename")
	if err != nil {
		log.Error("Error getting cookie", zap.Error(err))
	}
	filename = utils.RemoveSpecialCharacters(filename)
	database.DB = DB
	if !database.CheckIfTheHashHasCompleted(filename) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}
	filePath := utils.GetFilePath(filename)
	setHeaders(w, filename)
	http.ServeFile(w, r, filePath)
	utils.DeleteFile(filePath)
}

// Delete File
func DeleteFile(filename string) {
	err := utils.DeleteFile(filename)
	if err != nil {
		log.Error("Error deleting file after download", zap.Error(err))
	}
}
