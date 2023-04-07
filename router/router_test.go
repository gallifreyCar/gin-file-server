package router

import (
	"bytes"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

//todo Extract code to generate files.

func TestSetupRouter(t *testing.T) {

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
