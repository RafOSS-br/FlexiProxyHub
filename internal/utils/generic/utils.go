package utils

import (
	"crypto/md5"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"regexp"
	"go.uber.org/zap"
	"time"
	"strconv"
)

// Create file name from ci_session cookie, timestamp and path
func CreateFileName(r *http.Request) (string, error) {
	ci_session, err := getCISessionCookie(r)
	if err != nil {
		return "", err
	}
	return MD5(r.URL.Path + ci_session + r.URL.Path + GetTimestamp()), nil
}

// Create file with md5 hash to sinalize that file was downloaded
func CreateFileFlagWithHash(filename string) error {
	tempFolder := replaceBackslashWithSlash(getTempFolderPath())
	err := os.WriteFile(concatSlash(tempFolder)+filename+"flag", []byte("1"), 0644)
	if err != nil {
		return err
	}
	return nil
}

// check if file flag exists
func CheckIfFileFlagExists(filename string) bool {
	tempFolder := replaceBackslashWithSlash(getTempFolderPath())
	if _, err := os.Stat(concatSlash(tempFolder)+filename+"flag"); err == nil {
		return true
	}
	return false
}

// Get filepath

func GetFilePath(filename string) string {
	tempFolder := replaceBackslashWithSlash(getTempFolderPath())
	return concatSlash(tempFolder)+filename
}

// Get filename from cookie and remove special caracteres and file acess bypass
func GetFilenameFromCookie(r *http.Request) (string, error) {
	filename, err := GetCookie(r, "filename")
	if err != nil {
		return "", err
	}
	return RemoveSpecialCharacters(filename), nil
}

// Delete file flag
func DeleteFileFlag(filename string) error {
	tempFolder := replaceBackslashWithSlash(getTempFolderPath())
	err := os.Remove(concatSlash(tempFolder)+filename+"flag")
	if err != nil {
		return err
	}
	return nil
}

// Get timestamp
func GetTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// Get ci_session
func getCISessionCookie(r *http.Request) (string, error) {
    return GetCookie(r, "ci_session")
}

// Get cookie

func GetCookie (r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// MD5 hash from string
func MD5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return RemoveSpecialCharacters(Base64Encode(hash.Sum(nil)))

}

// remove special characters from string
func RemoveSpecialCharacters(text string) string {
    reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
    processedString := reg.ReplaceAllString(text, "")
    return processedString
} 

// base64 encode from []byte
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// save response body to file
func SaveResponseBodyToFile(filename string, recorder *httptest.ResponseRecorder, Log *zap.Logger) {
	body := recorder.Body.Bytes()
	tempFolder := replaceBackslashWithSlash(getTempFolderPath())
	err := os.WriteFile(concatSlash(tempFolder)+filename, body, 0644)
	if err != nil {
		Log.Error("Error writing file", zap.Error(err))
	}
}

// get temporary folder path windows and linux
func getTempFolderPath() string {
	return os.TempDir()
}

// replace \ with / in path
func replaceBackslashWithSlash(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

// if string last position not is / concat /
func concatSlash(path string) string {
	if path[len(path)-1:] != "/" {
		return path + "/"
	}
	return path
}

// convert string to int and return error if not possible
func ConvertStringToInt(text string) (int, error) {
	return strconv.Atoi(text)
}

//convert int to string
func ConvertIntToString(number int) string {
	return strconv.Itoa(number)
}

//delete file
func DeleteFile(filename string) error {
	return os.Remove(filename)
}