package database

import (
	"FlexiProxyHub/internal/database"
	"time"
)

type Database = database.Database
type FileStore = database.FileStore

var DB *database.Database

func Init() {
	DB = database.NewDatabase()
}

// Save Hash and Datetime to expiration
func Save(hash string, datetime string, contentType string, contentDisposition string) {
	DB.Set(hash, FileStore{Timestamp: datetime,
		ContentType:        contentType,
		ContentDisposition: contentDisposition})
}

// Get Datetime from Hash
func GetDatetime(hash string) string {
	return DB.Get(hash).Timestamp
}

// Delete Hash
func Delete(hash string) {
	DB.Delete(hash)
}

// Check if hash expired based in datetime timestamp based in now - 7 hours
func CheckIfExpired(timestamp string) bool {
	if timestamp == "" || timestamp == "downloading" {
		return true
	}
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return true
	}
	now := time.Now()
	if now.Sub(t).Hours() > 7 {
		return true
	} else {
		return false
	}
}

// Generate standard datetime timestamp

func GenerateTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// Save from hash and generate timestamp
func SaveFromHashAndComplete(hash string, contentType string, contentDisposition string) {
	Save(hash, GenerateTimestamp(), contentType, contentDisposition)
}

// Check if hash exists
func CheckIfHashExists(hash string) bool {
	if DB.Get(hash).Timestamp != "" {
		return true
	} else {
		return false
	}
}

// Check hash and return true if hash exists and not expired
func CheckIfTheHashHasCompleted(hash string) bool {
	if CheckIfHashExists(hash) {
		if !CheckIfExpired(DB.Get(hash).Timestamp) {
			return true
		}
	}
	return false
}

// Check if hash flagged as downloading
func CheckIfHashFlaggedAsDownloading(hash string) bool {
	return DB.Get(hash).Timestamp == "downloading"
}

// Get content type and content disposition from hash
func GetContentTypeAndContentDisposition(hash string) (string, string) {
	return DB.Get(hash).ContentType, DB.Get(hash).ContentDisposition
}

// Flag hash as downloading
func FlagHashAsDownloading(hash string) {
	DB.Set(hash, FileStore{Timestamp: "downloading",
		ContentType:        "",
		ContentDisposition: ""})
}

// Select all expired
func SelectAllExpired() []string {
	var expired []string
	for hash, fs := range DB.GetAll() {
		if CheckIfExpired(fs.Timestamp) {
			expired = append(expired, hash)
		}
	}
	return expired
}
