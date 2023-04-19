package router

import (
	"bytes"
	"crypto/rand"
	"github.com/gallifreyCar/gin-file-server/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestSetupRouter_old(t *testing.T) {

	testCases := []struct {
		name        string
		endpoint    string
		method      string
		contendType string
		expect      int
	}{
		{"TestUploadFileSingle", "/uploadFile", http.MethodPost, "multipart/form-data", http.StatusBadRequest},
		{"TestUploadFiles", "/uploadFiles", http.MethodPost, "multipart/form-data", http.StatusBadRequest},
		{"TestDownloadFile", "/downloadFile/single/myfile.txt", "", http.MethodGet, http.StatusNotFound},
		{"TestSelectFileLogByName", "/uploads/file_name/example3450268643.txt", "", http.MethodGet, http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setupRouter()
			w := httptest.NewRecorder()
			resp, err := http.NewRequest(tc.method, tc.endpoint, nil)
			resp.Header.Set("Content-Type", tc.contendType)
			router.ServeHTTP(w, resp)
			assert.Nil(t, err)
			assert.Equal(t, tc.expect, w.Code)
		})
	}

}

// Test MaxAllowed Middleware 9Mib
func TestRMMM(t *testing.T) {

	// set fileName and fileSize
	fileName := "testMMMFile"
	fileSize := int64(10 * 1024 * 1024) // 10 MB

	// generate len(data) random bytes and writes them into data
	data := make([]byte, fileSize)
	_, err := rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	// Create a temporary file to upload
	file, err := os.CreateTemp("", fileName+"*.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Write data into file
	if _, err := file.Write(data); err != nil {
		t.Fatal(err)
	}

	// Create a new request with the file as the body of the request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if _, err = io.Copy(part, file); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}

	// Create a test router
	router := setupRouter()
	// Create a response recorder to record the response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploadFile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("unexpected status code %d", w.Code)
		return
	}

	// Close file
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

}

func Test_fileUploadSingle(t *testing.T) {

	// Create a temporary file to upload
	files, clean, err := utils.CreateTempFiles(1, 1024*1024*2, "example*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	// Create a new request with the file as the body of the request
	body, writer, err := utils.CreateUploadBody("file", files)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test router
	router := setupRouter()
	// Create a response recorder to record the response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploadFile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	//Serve the request using the test router
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", w.Code)
		return
	}

}

func Test_fileUploadMultiple(t *testing.T) {
	fieldName := "testfiles"

	// Create temporary files to upload
	files, clean, err := utils.CreateTempFiles(2, 1024*1024*2, "example*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	// Create a new request with the file as the body of the request
	body, writer, err := utils.CreateUploadBody("file", files)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test router
	router := setupRouter()
	// Create a response recorder to record the response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploadFiles", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.PostForm = url.Values{"fieldName": {fieldName}}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	// Serve the request using the test router
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", w.Code)
		return
	}

}

// todo use testCase struct to rewrite this unit test
func TestDownloadFile(t *testing.T) {
	folder := "single"
	fileName := "testfile.txt"
	folder2 := "multiple"
	fileName2 := "testfile.txt"
	downloadAndSave(folder, fileName, t)
	downloadAndSave(folder2, fileName2, t)
}

func downloadAndSave(folder, fileName string, t *testing.T) {
	// Create a test router
	router := setupRouter()
	// Create a response recorder to record the response
	w := httptest.NewRecorder()
	requestUrl := path.Join("/downloadFile", folder, fileName)
	req, _ := http.NewRequest("GET", requestUrl, nil)

	// Serve the request using the test router
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code == http.StatusNotFound {
		t.Logf("file not found %d", w.Code)
		return
	}
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", w.Code)
		return
	}

	// Create the download directory if it doesn't exist
	downloadDirPath := path.Join("..", "target", "download", folder)
	// Save the files in the local
	saveFilePath := path.Join(downloadDirPath, fileName)

	err := os.MkdirAll(downloadDirPath, os.ModePerm)
	if err != nil {
		t.Errorf("error creating download directory: %v", err)
		return
	}

	// Create a new file to save the downloaded content
	out, err := os.Create(saveFilePath)
	if err != nil {
		t.Errorf("error creating file: %v", err)
		return
	}
	defer func() {
		err = out.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Write the downloaded content to the new file
	_, err = io.Copy(out, w.Body)
	if err != nil {
		t.Errorf("error writing file: %v", err)
		return
	}

}
